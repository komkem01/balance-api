package wallets

import (
	"balance/app/utils/base"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ListRequestService struct {
	MemberID *string `json:"member_id"`
	IsActive *bool   `json:"is_active"`
	Page     int     `json:"page"`
	Size     int     `json:"size"`
}

type ListItemService struct {
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

func (s *Service) ListWallet(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	items, err := s.db.ListWallets(ctx, req.IsActive)
	if err != nil {
		return nil, nil, err
	}
	res := make([]*ListItemService, 0, len(items))
	var memberFilter *uuid.UUID
	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				return nil, nil, ErrWalletInvalidMemberID
			}
			memberFilter = &id
		}
	}

	for _, item := range items {
		if memberFilter != nil {
			if item.MemberID == nil || *item.MemberID != *memberFilter {
				continue
			}
		}

		res = append(res, &ListItemService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Balance: item.Balance, Currency: item.Currency, ColorCode: item.ColorCode, IsActive: item.IsActive, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt})
	}

	page := int64(req.Page)
	if page < 1 {
		page = 1
	}
	size := int64(req.Size)
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	total := int64(len(res))
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}
	return res[start:end], &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
