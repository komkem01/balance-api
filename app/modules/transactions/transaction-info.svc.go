package transactions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string `json:"id"`
}

type InfoResponseService struct {
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

func (s *Service) InfoTransaction(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrTransactionInvalidID
	}
	item, err := s.db.GetTransactionByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return &InfoResponseService{ID: item.ID, WalletID: item.WalletID, CategoryID: item.CategoryID, Amount: item.Amount, Type: item.Type, TransactionDate: item.TransactionDate, Note: item.Note, ImageURL: item.ImageURL, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
