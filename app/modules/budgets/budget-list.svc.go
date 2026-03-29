package budgets

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
	CategoryID *string `json:"category_id"`
	Period     *string `json:"period"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
}

type ListItemService struct {
	ID         uuid.UUID        `json:"id"`
	MemberID   *uuid.UUID       `json:"member_id"`
	CategoryID *uuid.UUID       `json:"category_id"`
	Amount     float64          `json:"amount"`
	Period     ent.BudgetPeriod `json:"period"`
	StartDate  *time.Time       `json:"start_date"`
	EndDate    *time.Time       `json:"end_date"`
	CreatedAt  time.Time        `json:"created_at"`
}

func (s *Service) ListBudget(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	var period *ent.BudgetPeriod
	if req.Period != nil {
		v := strings.TrimSpace(*req.Period)
		if v != "" {
			parsed, ok := parseBudgetPeriod(v)
			if !ok {
				return nil, nil, ErrBudgetPeriodInvalid
			}
			period = &parsed
		}
	}

	items, err := s.db.ListBudgets(ctx, req.MemberID, req.CategoryID, period)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{ID: item.ID, MemberID: item.MemberID, CategoryID: item.CategoryID, Amount: item.Amount, Period: item.Period, StartDate: item.StartDate, EndDate: item.EndDate, CreatedAt: item.CreatedAt})
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
