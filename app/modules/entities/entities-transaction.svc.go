package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.TransactionEntity = (*Service)(nil)

var ErrWalletBalanceInsufficient = errors.New("wallet balance would be negative")

func parseTransactionID(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func transactionWalletDelta(amount float64, transactionType ent.TransactionType) float64 {
	if transactionType == ent.TransactionTypeExpense {
		return -amount
	}
	return amount
}

func (s *Service) adjustWalletBalanceTx(ctx context.Context, tx bun.Tx, walletID uuid.UUID, delta float64) error {
	wallet := &ent.WalletEntity{}
	if err := tx.NewSelect().
		Model(wallet).
		Where("wallet.id = ?", walletID).
		Where("wallet.deleted_at IS NULL").
		For("UPDATE").
		Scan(ctx); err != nil {
		return err
	}

	newBalance := wallet.Balance + delta
	if newBalance < 0 {
		return ErrWalletBalanceInsufficient
	}
	wallet.Balance = newBalance
	_, err := tx.NewUpdate().
		Model(wallet).
		WherePK().
		Column("balance", "updated_at").
		Exec(ctx)
	return err
}

func (s *Service) CreateTransaction(ctx context.Context, walletID *string, categoryID *string, amount float64, transactionType ent.TransactionType, transactionDate *time.Time, note string, imageURL string) (*ent.TransactionEntity, error) {
	wid, err := parseTransactionID(walletID)
	if err != nil {
		return nil, err
	}
	cid, err := parseTransactionID(categoryID)
	if err != nil {
		return nil, err
	}

	model := &ent.TransactionEntity{
		ID:              uuid.New(),
		WalletID:        wid,
		CategoryID:      cid,
		Amount:          amount,
		Type:            transactionType,
		TransactionDate: transactionDate,
		Note:            strings.TrimSpace(note),
		ImageURL:        strings.TrimSpace(imageURL),
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) CreateTransactionWithWalletAdjust(ctx context.Context, walletID *string, categoryID *string, amount float64, transactionType ent.TransactionType, transactionDate *time.Time, note string, imageURL string) (*ent.TransactionEntity, error) {
	wid, err := parseTransactionID(walletID)
	if err != nil {
		return nil, err
	}
	cid, err := parseTransactionID(categoryID)
	if err != nil {
		return nil, err
	}

	model := &ent.TransactionEntity{
		ID:              uuid.New(),
		WalletID:        wid,
		CategoryID:      cid,
		Amount:          amount,
		Type:            transactionType,
		TransactionDate: transactionDate,
		Note:            strings.TrimSpace(note),
		ImageURL:        strings.TrimSpace(imageURL),
	}

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if model.WalletID != nil {
			delta := transactionWalletDelta(model.Amount, model.Type)
			if err := s.adjustWalletBalanceTx(ctx, tx, *model.WalletID, delta); err != nil {
				return err
			}
		}

		if _, err := tx.NewInsert().Model(model).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id string) (*ent.TransactionEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.TransactionEntity{}
	if err := s.db.NewSelect().Model(model).Where("transaction.id = ?", uid).Where("transaction.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateTransaction(ctx context.Context, id string, walletID *string, categoryID *string, amount *float64, transactionType *ent.TransactionType, transactionDate *time.Time, note *string, imageURL *string) (*ent.TransactionEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.TransactionEntity{}
	if err := s.db.NewSelect().Model(model).Where("transaction.id = ?", uid).Scan(ctx); err != nil {
		return nil, err
	}

	if walletID != nil {
		wid, err := parseTransactionID(walletID)
		if err != nil {
			return nil, err
		}
		model.WalletID = wid
	}
	if categoryID != nil {
		cid, err := parseTransactionID(categoryID)
		if err != nil {
			return nil, err
		}
		model.CategoryID = cid
	}
	if amount != nil {
		model.Amount = *amount
	}
	if transactionType != nil {
		model.Type = *transactionType
	}
	if transactionDate != nil {
		model.TransactionDate = transactionDate
	}
	if note != nil {
		model.Note = strings.TrimSpace(*note)
	}
	if imageURL != nil {
		model.ImageURL = strings.TrimSpace(*imageURL)
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("wallet_id", "category_id", "amount", "type", "transaction_date", "note", "image_url", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateTransactionWithWalletAdjust(ctx context.Context, id string, walletID *string, categoryID *string, amount *float64, transactionType *ent.TransactionType, transactionDate *time.Time, note *string, imageURL *string) (*ent.TransactionEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var updated *ent.TransactionEntity
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		model := &ent.TransactionEntity{}
		if err := tx.NewSelect().
			Model(model).
			Where("transaction.id = ?", uid).
			Where("transaction.deleted_at IS NULL").
			For("UPDATE").
			Scan(ctx); err != nil {
			return err
		}

		oldWalletID := model.WalletID
		oldAmount := model.Amount
		oldType := model.Type

		if walletID != nil {
			wid, err := parseTransactionID(walletID)
			if err != nil {
				return err
			}
			model.WalletID = wid
		}
		if categoryID != nil {
			cid, err := parseTransactionID(categoryID)
			if err != nil {
				return err
			}
			model.CategoryID = cid
		}
		if amount != nil {
			model.Amount = *amount
		}
		if transactionType != nil {
			model.Type = *transactionType
		}
		if transactionDate != nil {
			model.TransactionDate = transactionDate
		}
		if note != nil {
			model.Note = strings.TrimSpace(*note)
		}
		if imageURL != nil {
			model.ImageURL = strings.TrimSpace(*imageURL)
		}

		oldDelta := transactionWalletDelta(oldAmount, oldType)
		newDelta := transactionWalletDelta(model.Amount, model.Type)

		if oldWalletID != nil && model.WalletID != nil && *oldWalletID == *model.WalletID {
			net := newDelta - oldDelta
			if net != 0 {
				if err := s.adjustWalletBalanceTx(ctx, tx, *model.WalletID, net); err != nil {
					return err
				}
			}
		} else {
			if oldWalletID != nil {
				if err := s.adjustWalletBalanceTx(ctx, tx, *oldWalletID, -oldDelta); err != nil {
					return err
				}
			}
			if model.WalletID != nil {
				if err := s.adjustWalletBalanceTx(ctx, tx, *model.WalletID, newDelta); err != nil {
					return err
				}
			}
		}

		if _, err := tx.NewUpdate().
			Model(model).
			WherePK().
			Column("wallet_id", "category_id", "amount", "type", "transaction_date", "note", "image_url", "updated_at").
			Exec(ctx); err != nil {
			return err
		}

		updated = model
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.db.NewUpdate().
		Model(&ent.TransactionEntity{}).
		Set("deleted_at = now()").
		Set("updated_at = now()").
		Where("id = ?", uid).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (s *Service) DeleteTransactionWithWalletAdjust(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		model := &ent.TransactionEntity{}
		if err := tx.NewSelect().
			Model(model).
			Where("transaction.id = ?", uid).
			Where("transaction.deleted_at IS NULL").
			For("UPDATE").
			Scan(ctx); err != nil {
			return err
		}

		if model.WalletID != nil {
			delta := -transactionWalletDelta(model.Amount, model.Type)
			if err := s.adjustWalletBalanceTx(ctx, tx, *model.WalletID, delta); err != nil {
				return err
			}
		}

		if _, err := tx.NewUpdate().
			Model(&ent.TransactionEntity{}).
			Set("deleted_at = now()").
			Set("updated_at = now()").
			Where("id = ?", uid).
			Where("deleted_at IS NULL").
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ListTransactions(ctx context.Context, walletID *string, categoryID *string, transactionType *ent.TransactionType) ([]*ent.TransactionEntity, error) {
	items := make([]*ent.TransactionEntity, 0)
	q := s.db.NewSelect().Model(&items).Where("transaction.deleted_at IS NULL").Order("transaction.created_at DESC")

	if walletID != nil {
		wid, err := parseTransactionID(walletID)
		if err != nil {
			return nil, err
		}
		if wid == nil {
			q = q.Where("transaction.wallet_id IS NULL")
		} else {
			q = q.Where("transaction.wallet_id = ?", *wid)
		}
	}

	if categoryID != nil {
		cid, err := parseTransactionID(categoryID)
		if err != nil {
			return nil, err
		}
		if cid == nil {
			q = q.Where("transaction.category_id IS NULL")
		} else {
			q = q.Where("transaction.category_id = ?", *cid)
		}
	}

	if transactionType != nil {
		q = q.Where("transaction.type = ?", *transactionType)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
