package transactions

import (
	"context"
	"strings"
	"time"

	"balance/app/modules/entities/ent"
	"balance/app/utils/base"

	"github.com/google/uuid"
)

type ListRequestService struct {
	MemberID   *string `json:"member_id"`
	WalletID   *string `json:"wallet_id"`
	CategoryID *string `json:"category_id"`
	Type       *string `json:"type"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
}

type ListItemService struct {
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

func (s *Service) ListTransaction(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	var transactionType *ent.TransactionType
	if req.Type != nil {
		v := strings.TrimSpace(*req.Type)
		if v != "" {
			parsed, ok := parseTransactionType(v)
			if !ok {
				return nil, nil, ErrTransactionTypeInvalid
			}
			transactionType = &parsed
		}
	}

	items, err := s.db.ListTransactions(ctx, req.MemberID, req.WalletID, req.CategoryID, transactionType)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{ID: item.ID, WalletID: item.WalletID, CategoryID: item.CategoryID, Amount: item.Amount, Type: item.Type, TransactionDate: item.TransactionDate, Note: item.Note, ImageURL: item.ImageURL, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt})
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
