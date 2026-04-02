package budgets

import (
	"context"
	"math"
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
	ID          uuid.UUID        `json:"id"`
	MemberID    *uuid.UUID       `json:"member_id"`
	CategoryID  *uuid.UUID       `json:"category_id"`
	Amount      float64          `json:"amount"`
	SpentAmount float64          `json:"spent_amount"`
	UsedPercent float64          `json:"used_percent"`
	Period      ent.BudgetPeriod `json:"period"`
	StartDate   *time.Time       `json:"start_date"`
	EndDate     *time.Time       `json:"end_date"`
	CreatedAt   time.Time        `json:"created_at"`
}

type ListResponseService struct {
	Items         []*ListItemService `json:"items"`
	TotalNetWorth float64            `json:"total_net_worth"`
}

func (s *Service) ListBudget(ctx context.Context, req *ListRequestService) (*ListResponseService, *base.ResponsePaginate, error) {
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

	walletByID := make(map[uuid.UUID]struct{})
	totalNetWorth := 0.0
	if req.MemberID != nil {
		isActive := true
		wallets, err := s.db.ListWallets(ctx, &isActive)
		if err != nil {
			return nil, nil, err
		}
		memberIDStr := strings.TrimSpace(*req.MemberID)
		for _, wallet := range wallets {
			if wallet.MemberID == nil {
				continue
			}
			if wallet.MemberID.String() != memberIDStr {
				continue
			}
			walletByID[wallet.ID] = struct{}{}
			totalNetWorth += wallet.Balance
		}
	}

	res := make([]*ListItemService, 0, len(items))
	loc := budgetLocation()
	for _, item := range items {
		effectiveStartDate := item.StartDate
		effectiveEndDate := item.EndDate
		if effectiveStartDate == nil || effectiveEndDate == nil {
			resolvedStart, resolvedEnd, err := resolveBudgetDateRange(item.Period, item.StartDate, item.EndDate)
			if err != nil {
				return nil, nil, err
			}
			effectiveStartDate = resolvedStart
			effectiveEndDate = resolvedEnd

			if _, err := s.db.UpdateBudget(ctx, item.ID.String(), nil, nil, nil, nil, effectiveStartDate, effectiveEndDate); err != nil {
				return nil, nil, err
			}
		}

		spentAmount := 0.0
		if len(walletByID) > 0 {
			cat := ""
			if item.CategoryID != nil {
				cat = item.CategoryID.String()
			}
			expenseType := ent.TransactionTypeExpense
			for walletID := range walletByID {
				walletIDStr := walletID.String()
				var categoryID *string
				if cat != "" {
					categoryID = &cat
				}
				txItems, err := s.db.ListTransactions(ctx, nil, &walletIDStr, categoryID, &expenseType)
				if err != nil {
					return nil, nil, err
				}
				for _, tx := range txItems {
					txTime := tx.CreatedAt
					if tx.TransactionDate != nil {
						txTime = *tx.TransactionDate
					}
					txDate := dateOnlyInLocation(txTime, loc)
					if effectiveStartDate != nil && txDate.Before(dateOnlyInLocation(*effectiveStartDate, loc)) {
						continue
					}
					if effectiveEndDate != nil && txDate.After(dateOnlyInLocation(*effectiveEndDate, loc)) {
						continue
					}
					spentAmount += tx.Amount
				}
			}
		}

		spentAmount = math.Round(spentAmount*100) / 100
		if math.Abs(item.SpentAmount-spentAmount) > 0.0001 {
			if err := s.db.UpdateBudgetSpent(ctx, item.ID.String(), spentAmount); err != nil {
				return nil, nil, err
			}
		}

		usedPercent := 0.0
		if item.Amount > 0 {
			usedPercent = math.Round((spentAmount/item.Amount)*10000) / 100
		}

		res = append(res, &ListItemService{
			ID:          item.ID,
			MemberID:    item.MemberID,
			CategoryID:  item.CategoryID,
			Amount:      item.Amount,
			SpentAmount: spentAmount,
			UsedPercent: usedPercent,
			Period:      item.Period,
			StartDate:   effectiveStartDate,
			EndDate:     effectiveEndDate,
			CreatedAt:   item.CreatedAt,
		})
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

	return &ListResponseService{Items: res[start:end], TotalNetWorth: totalNetWorth}, &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
