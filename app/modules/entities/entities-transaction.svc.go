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

func (s *Service) CreateTransferTransactionsWithWalletAdjust(ctx context.Context, fromWalletID string, toWalletID string, categoryID *string, amount float64, transactionDate *time.Time, fromNote string, toNote string) (*ent.TransactionEntity, *ent.TransactionEntity, error) {
	fromWalletUUID, err := uuid.Parse(strings.TrimSpace(fromWalletID))
	if err != nil {
		return nil, nil, err
	}

	toWalletUUID, err := uuid.Parse(strings.TrimSpace(toWalletID))
	if err != nil {
		return nil, nil, err
	}

	categoryUUID, err := parseTransactionID(categoryID)
	if err != nil {
		return nil, nil, err
	}

	fromNote = strings.TrimSpace(fromNote)
	if fromNote == "" {
		fromNote = "Wallet transfer"
	}

	toNote = strings.TrimSpace(toNote)
	if toNote == "" {
		toNote = "Wallet transfer"
	}

	fromTx := &ent.TransactionEntity{
		ID:              uuid.New(),
		WalletID:        &fromWalletUUID,
		CategoryID:      categoryUUID,
		Amount:          amount,
		Type:            ent.TransactionTypeExpense,
		TransactionDate: transactionDate,
		Note:            fromNote,
		ImageURL:        "",
	}

	toTx := &ent.TransactionEntity{
		ID:              uuid.New(),
		WalletID:        &toWalletUUID,
		CategoryID:      categoryUUID,
		Amount:          amount,
		Type:            ent.TransactionTypeIncome,
		TransactionDate: transactionDate,
		Note:            toNote,
		ImageURL:        "",
	}

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := s.adjustWalletBalanceTx(ctx, tx, fromWalletUUID, -amount); err != nil {
			return err
		}
		if err := s.adjustWalletBalanceTx(ctx, tx, toWalletUUID, amount); err != nil {
			return err
		}

		if _, err := tx.NewInsert().Model(fromTx).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(toTx).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return fromTx, toTx, nil
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

func (s *Service) ListTransactions(ctx context.Context, memberID *string, walletID *string, categoryID *string, transactionType *ent.TransactionType) ([]*ent.TransactionEntity, error) {
	items := make([]*ent.TransactionEntity, 0)
	q := s.db.NewSelect().
		Model(&items).
		Join("JOIN wallets AS wallet ON wallet.id = transaction.wallet_id").
		Where("transaction.deleted_at IS NULL").
		Where("wallet.deleted_at IS NULL").
		Order("transaction.created_at DESC")

	if memberID != nil {
		mid, err := parseTransactionID(memberID)
		if err != nil {
			return nil, err
		}
		if mid == nil {
			q = q.Where("1 = 0")
		} else {
			q = q.Where("wallet.member_id = ?", *mid)
		}
	}

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

func (s *Service) ListTransactionMonthlySummary(ctx context.Context, memberID *string, walletID *string, categoryID *string, startDate *time.Time, endDate *time.Time) ([]*ent.TransactionMonthlySummaryEntity, error) {
	items := make([]*ent.TransactionMonthlySummaryEntity, 0)

	q := s.db.NewSelect().
		TableExpr("transactions AS transaction").
		ColumnExpr("date_trunc('month', COALESCE(transaction.transaction_date::timestamp, transaction.created_at))::date AS month").
		ColumnExpr("COALESCE(SUM(CASE WHEN transaction.type = ? THEN transaction.amount ELSE 0 END), 0) AS income_total", ent.TransactionTypeIncome).
		ColumnExpr("COALESCE(SUM(CASE WHEN transaction.type = ? THEN transaction.amount ELSE 0 END), 0) AS expense_total", ent.TransactionTypeExpense).
		ColumnExpr("COUNT(*) AS transaction_count").
		Join("JOIN wallets AS wallet ON wallet.id = transaction.wallet_id").
		Where("transaction.deleted_at IS NULL").
		Where("transaction.note NOT LIKE ?", "__transfer__|%").
		Where("transaction.note <> ?", "Wallet transfer").
		Where("wallet.deleted_at IS NULL").
		GroupExpr("1").
		OrderExpr("1 ASC")

	if memberID != nil {
		mid, err := parseTransactionID(memberID)
		if err != nil {
			return nil, err
		}
		if mid == nil {
			q = q.Where("1 = 0")
		} else {
			q = q.Where("wallet.member_id = ?", *mid)
		}
	}

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

	if startDate != nil {
		q = q.Where("COALESCE(transaction.transaction_date::timestamp, transaction.created_at) >= ?", *startDate)
	}

	if endDate != nil {
		nextDay := endDate.AddDate(0, 0, 1)
		q = q.Where("COALESCE(transaction.transaction_date::timestamp, transaction.created_at) < ?", nextDay)
	}

	if err := q.Scan(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}
