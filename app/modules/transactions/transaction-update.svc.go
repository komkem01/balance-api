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

type UpdateRequestService struct {
	ID              string   `json:"id"`
	WalletID        *string  `json:"wallet_id"`
	CategoryID      *string  `json:"category_id"`
	Amount          *float64 `json:"amount"`
	Type            *string  `json:"type"`
	TransactionDate *string  `json:"transaction_date"`
	Note            *string  `json:"note"`
	ImageURL        *string  `json:"image_url"`
}

type UpdateResponseService struct {
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

func (s *Service) UpdateTransaction(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrTransactionInvalidID
	}
	if req.WalletID == nil && req.CategoryID == nil && req.Amount == nil && req.Type == nil && req.TransactionDate == nil && req.Note == nil && req.ImageURL == nil {
		return nil, ErrTransactionNoFieldsToUpdate
	}
	if req.Amount != nil && *req.Amount < 0 {
		return nil, ErrTransactionAmountInvalid
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

	var transactionType *ent.TransactionType
	if req.Type != nil {
		parsed, ok := parseTransactionType(strings.TrimSpace(*req.Type))
		if !ok {
			return nil, ErrTransactionTypeInvalid
		}
		transactionType = &parsed
	}

	transactionDate, err := parseDateString(req.TransactionDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}

	item, err := s.db.UpdateTransactionWithWalletAdjust(ctx, req.ID, req.WalletID, req.CategoryID, req.Amount, transactionType, transactionDate, req.Note, req.ImageURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		if errors.Is(err, entities.ErrWalletBalanceInsufficient) {
			return nil, ErrTransactionInsufficientFunds
		}
		return nil, err
	}

	return &UpdateResponseService{ID: item.ID, WalletID: item.WalletID, CategoryID: item.CategoryID, Amount: item.Amount, Type: item.Type, TransactionDate: item.TransactionDate, Note: item.Note, ImageURL: s.resolveImageURL(ctx, item.ImageURL), CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
