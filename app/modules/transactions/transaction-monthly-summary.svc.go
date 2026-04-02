package transactions

import (
	"context"
	"strings"
	"time"
)

type MonthlySummaryRequestService struct {
	MemberID   *string `json:"member_id"`
	WalletID   *string `json:"wallet_id"`
	CategoryID *string `json:"category_id"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
}

type MonthlySummaryItemService struct {
	Month            string  `json:"month"`
	IncomeTotal      float64 `json:"income_total"`
	ExpenseTotal     float64 `json:"expense_total"`
	TransactionCount int64   `json:"transaction_count"`
}

func parseDateOnly(value *string) (*time.Time, error) {
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

func (s *Service) MonthlySummaryTransaction(ctx context.Context, req *MonthlySummaryRequestService) ([]*MonthlySummaryItemService, error) {
	startDate, err := parseDateOnly(req.StartDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}
	endDate, err := parseDateOnly(req.EndDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, ErrTransactionDateInvalid
	}

	items, err := s.db.ListTransactionMonthlySummary(ctx, req.MemberID, req.WalletID, req.CategoryID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	res := make([]*MonthlySummaryItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &MonthlySummaryItemService{
			Month:            item.Month.Format("2006-01"),
			IncomeTotal:      item.IncomeTotal,
			ExpenseTotal:     item.ExpenseTotal,
			TransactionCount: item.TransactionCount,
		})
	}

	return res, nil
}
