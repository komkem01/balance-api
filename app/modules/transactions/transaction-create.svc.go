package transactions

import (
	"balance/app/modules/entities"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	WalletID        *string `json:"wallet_id"`
	CategoryID      *string `json:"category_id"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	TransactionDate *string `json:"transaction_date"`
	Note            string  `json:"note"`
	ImageURL        string  `json:"image_url"`
}

type CreateResponseService struct {
	ID              uuid.UUID           `json:"id"`
	WalletID        *uuid.UUID          `json:"wallet_id"`
	CategoryID      *uuid.UUID          `json:"category_id"`
	Amount          float64             `json:"amount"`
	Type            ent.TransactionType `json:"type"`
	TransactionDate *time.Time          `json:"transaction_date"`
	Note            string              `json:"note"`
	ImageURL        string              `json:"image_url"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

func parseDateString(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Service) CreateTransaction(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if req.Amount < 0 {
		return nil, ErrTransactionAmountInvalid
	}

	transactionType, ok := parseTransactionType(strings.TrimSpace(req.Type))
	if !ok {
		return nil, ErrTransactionTypeInvalid
	}

	if req.WalletID != nil {
		v := strings.TrimSpace(*req.WalletID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrTransactionInvalidWalletID
			}
			if _, err := s.db.GetWalletByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrTransactionInvalidWalletID
				}
				return nil, err
			}
		}
	}

	if req.CategoryID != nil {
		v := strings.TrimSpace(*req.CategoryID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrTransactionInvalidCategoryID
			}
			if _, err := s.db.GetCategoryByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrTransactionInvalidCategoryID
				}
				return nil, err
			}
		}
	}

	transactionDate, err := parseDateString(req.TransactionDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}

	item, err := s.db.CreateTransactionWithWalletAdjust(
		ctx,
		req.WalletID,
		req.CategoryID,
		req.Amount,
		transactionType,
		transactionDate,
		strings.TrimSpace(req.Note),
		strings.TrimSpace(req.ImageURL),
	)
	if err != nil {
		if errors.Is(err, entities.ErrWalletBalanceInsufficient) {
			return nil, ErrTransactionInsufficientFunds
		}
		return nil, err
	}

	return &CreateResponseService{
		ID:              item.ID,
		WalletID:        item.WalletID,
		CategoryID:      item.CategoryID,
		Amount:          item.Amount,
		Type:            item.Type,
		TransactionDate: item.TransactionDate,
		Note:            item.Note,
		ImageURL:        item.ImageURL,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}, nil
}
