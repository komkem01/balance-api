package loans

import (
	"context"
	"time"

	"balance/app/utils/base"

	"github.com/google/uuid"
)

type ListRequestService struct {
	MemberID *string `json:"member_id"`
	Page     int     `json:"page"`
	Size     int     `json:"size"`
}

type ListItemService struct {
	ID               uuid.UUID  `json:"id"`
	MemberID         *uuid.UUID `json:"member_id"`
	Name             string     `json:"name"`
	Lender           string     `json:"lender"`
	TotalAmount      float64    `json:"total_amount"`
	RemainingBalance float64    `json:"remaining_balance"`
	MonthlyPayment   float64    `json:"monthly_payment"`
	InterestRate     float64    `json:"interest_rate"`
	ColorCode        string     `json:"color_code"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (s *Service) ListLoan(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	items, err := s.db.ListLoans(ctx, req.MemberID)
	if err != nil {
		return nil, nil, err
	}
	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{
			ID:               item.ID,
			MemberID:         item.MemberID,
			Name:             item.Name,
			Lender:           item.Lender,
			TotalAmount:      item.TotalAmount,
			RemainingBalance: item.RemainingBalance,
			MonthlyPayment:   item.MonthlyPayment,
			InterestRate:     item.InterestRate,
			ColorCode:        item.ColorCode,
			StartDate:        item.StartDate,
			EndDate:          item.EndDate,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		})
	}

	total := int64(len(items))
	page := int64(req.Page)
	size := int64(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	return res, &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
