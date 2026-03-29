package wallets

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string `json:"id"`
}

type InfoResponseService struct {
	ID        uuid.UUID  `json:"id"`
	MemberID  *uuid.UUID `json:"member_id"`
	Name      string     `json:"name"`
	Balance   float64    `json:"balance"`
	Currency  string     `json:"currency"`
	ColorCode string     `json:"color_code"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Service) InfoWallet(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrWalletInvalidID
	}
	item, err := s.db.GetWalletByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return &InfoResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Balance: item.Balance, Currency: item.Currency, ColorCode: item.ColorCode, IsActive: item.IsActive, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
