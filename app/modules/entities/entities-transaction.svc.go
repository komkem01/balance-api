package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.TransactionEntity = (*Service)(nil)

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

func (s *Service) GetTransactionByID(ctx context.Context, id string) (*ent.TransactionEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.TransactionEntity{}
	if err := s.db.NewSelect().Model(model).Where("transaction.id = ?", uid).Scan(ctx); err != nil {
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

func (s *Service) DeleteTransaction(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.db.NewDelete().Model(&ent.TransactionEntity{}).Where("id = ?", uid).Exec(ctx)
	return err
}

func (s *Service) ListTransactions(ctx context.Context, walletID *string, categoryID *string, transactionType *ent.TransactionType) ([]*ent.TransactionEntity, error) {
	items := make([]*ent.TransactionEntity, 0)
	q := s.db.NewSelect().Model(&items).Order("transaction.created_at DESC")

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
