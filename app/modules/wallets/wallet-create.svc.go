package wallets

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	MemberID  *string `json:"member_id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	ColorCode string  `json:"color_code"`
	IsActive  *bool   `json:"is_active"`
}

type CreateResponseService struct {
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

func (s *Service) CreateWallet(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrWalletNameRequired
	}

	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrWalletInvalidMemberID
			}
			if _, err := s.db.GetMemberByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrWalletInvalidMemberID
				}
				return nil, err
			}
		}
	}

	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		currency = "THB"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := s.db.CreateWallet(ctx, req.MemberID, strings.TrimSpace(req.Name), req.Balance, currency, strings.TrimSpace(req.ColorCode), isActive)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Balance: item.Balance, Currency: item.Currency, ColorCode: item.ColorCode, IsActive: item.IsActive, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
