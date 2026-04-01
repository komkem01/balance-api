package budgets

import (
	"context"
	"math"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type RecalculateAllResponseService struct {
	TotalBudgets       int       `json:"total_budgets"`
	UpdatedDateRanges  int       `json:"updated_date_ranges"`
	UpdatedSpentAmount int       `json:"updated_spent_amount"`
	RecalculatedAt     time.Time `json:"recalculated_at"`
}

func sameCalendarDay(a *time.Time, b *time.Time, loc *time.Location) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ad := dateOnlyInLocation(*a, loc)
	bd := dateOnlyInLocation(*b, loc)
	return ad.Equal(bd)
}

func (s *Service) RecalculateAllBudgets(ctx context.Context) (*RecalculateAllResponseService, error) {
	items, err := s.db.ListBudgets(ctx, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	wallets, err := s.db.ListWallets(ctx, nil)
	if err != nil {
		return nil, err
	}

	walletIDsByMember := make(map[uuid.UUID][]string)
	for _, wallet := range wallets {
		if wallet.MemberID == nil {
			continue
		}
		walletIDsByMember[*wallet.MemberID] = append(walletIDsByMember[*wallet.MemberID], wallet.ID.String())
	}

	loc := budgetLocation()
	expenseType := ent.TransactionTypeExpense
	updatedRanges := 0
	updatedSpent := 0

	for _, item := range items {
		startDate, endDate, err := resolveBudgetDateRange(item.Period, item.StartDate, item.EndDate)
		if err != nil {
			return nil, err
		}

		if !sameCalendarDay(item.StartDate, startDate, loc) || !sameCalendarDay(item.EndDate, endDate, loc) {
			if _, err := s.db.UpdateBudget(ctx, item.ID.String(), nil, nil, nil, nil, startDate, endDate); err != nil {
				return nil, err
			}
			updatedRanges++
		}

		spentAmount := 0.0
		if item.MemberID != nil {
			walletIDs := walletIDsByMember[*item.MemberID]
			category := ""
			if item.CategoryID != nil {
				category = item.CategoryID.String()
			}

			for _, walletID := range walletIDs {
				var categoryID *string
				if category != "" {
					categoryID = &category
				}
				txItems, err := s.db.ListTransactions(ctx, &walletID, categoryID, &expenseType)
				if err != nil {
					return nil, err
				}

				for _, tx := range txItems {
					txTime := tx.CreatedAt
					if tx.TransactionDate != nil {
						txTime = *tx.TransactionDate
					}
					txDate := dateOnlyInLocation(txTime, loc)
					if startDate != nil && txDate.Before(dateOnlyInLocation(*startDate, loc)) {
						continue
					}
					if endDate != nil && txDate.After(dateOnlyInLocation(*endDate, loc)) {
						continue
					}
					spentAmount += tx.Amount
				}
			}
		}

		spentAmount = math.Round(spentAmount*100) / 100
		if math.Abs(item.SpentAmount-spentAmount) > 0.0001 {
			if err := s.db.UpdateBudgetSpent(ctx, item.ID.String(), spentAmount); err != nil {
				return nil, err
			}
			updatedSpent++
		}
	}

	return &RecalculateAllResponseService{
		TotalBudgets:       len(items),
		UpdatedDateRanges:  updatedRanges,
		UpdatedSpentAmount: updatedSpent,
		RecalculatedAt:     time.Now().UTC(),
	}, nil
}
