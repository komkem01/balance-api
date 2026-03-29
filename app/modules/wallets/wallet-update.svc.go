package wallets

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID        string   `json:"id"`
	MemberID  *string  `json:"member_id"`
	Name      *string  `json:"name"`
	Balance   *float64 `json:"balance"`
	Currency  *string  `json:"currency"`
	ColorCode *string  `json:"color_code"`
	IsActive  *bool    `json:"is_active"`
}

type UpdateResponseService struct {
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

func (s *Service) UpdateWallet(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrWalletInvalidID
	}
	if req.MemberID == nil && req.Name == nil && req.Balance == nil && req.Currency == nil && req.ColorCode == nil && req.IsActive == nil {
		return nil, ErrWalletNoFieldsToUpdate
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
	item, err := s.db.UpdateWallet(ctx, req.ID, req.MemberID, req.Name, req.Balance, req.Currency, req.ColorCode, req.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return &UpdateResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Balance: item.Balance, Currency: item.Currency, ColorCode: item.ColorCode, IsActive: item.IsActive, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}, nil
}
