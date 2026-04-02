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
	Range      *string `json:"range"`
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

func parseSummaryRange(value *string) (*time.Time, *time.Time, error) {
	if value == nil {
		return nil, nil, nil
	}

	v := strings.ToLower(strings.TrimSpace(*value))
	if v == "" || v == "all" {
		return nil, nil, nil
	}

	now := time.Now()
	switch v {
	case "1d":
		return ptrTime(now.Add(-24 * time.Hour)), ptrTime(now), nil
	case "1w":
		return ptrTime(now.AddDate(0, 0, -7)), ptrTime(now), nil
	case "1m":
		return ptrTime(now.AddDate(0, -1, 0)), ptrTime(now), nil
	case "1y":
		return ptrTime(now.AddDate(-1, 0, 0)), ptrTime(now), nil
	default:
		return nil, nil, ErrTransactionRangeInvalid
	}
}

func ptrTime(value time.Time) *time.Time {
	v := value
	return &v
}

func (s *Service) MonthlySummaryTransaction(ctx context.Context, req *MonthlySummaryRequestService) ([]*MonthlySummaryItemService, error) {
	startDate, endDate, err := parseSummaryRange(req.Range)
	if err != nil {
		return nil, err
	}

	explicitStartDate, err := parseDateOnly(req.StartDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}
	explicitEndDate, err := parseDateOnly(req.EndDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}

	if explicitStartDate != nil {
		startDate = explicitStartDate
	}
	if explicitEndDate != nil {
		endDate = explicitEndDate
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
