package budgets

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID         string   `json:"id"`
	MemberID   *string  `json:"member_id"`
	CategoryID *string  `json:"category_id"`
	Amount     *float64 `json:"amount"`
	Period     *string  `json:"period"`
	StartDate  *string  `json:"start_date"`
	EndDate    *string  `json:"end_date"`
}

type UpdateResponseService struct {
	ID          uuid.UUID        `json:"id"`
	MemberID    *uuid.UUID       `json:"member_id"`
	CategoryID  *uuid.UUID       `json:"category_id"`
	Amount      float64          `json:"amount"`
	SpentAmount float64          `json:"spent_amount"`
	Period      ent.BudgetPeriod `json:"period"`
	StartDate   *time.Time       `json:"start_date"`
	EndDate     *time.Time       `json:"end_date"`
	CreatedAt   time.Time        `json:"created_at"`
}

func (s *Service) UpdateBudget(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrBudgetInvalidID
	}
	if req.MemberID == nil && req.CategoryID == nil && req.Amount == nil && req.Period == nil && req.StartDate == nil && req.EndDate == nil {
		return nil, ErrBudgetNoFieldsToUpdate
	}

	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrBudgetInvalidMemberID
			}
			if _, err := s.db.GetMemberByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrBudgetInvalidMemberID
				}
				return nil, err
			}
		}
	}

	if req.CategoryID != nil {
		v := strings.TrimSpace(*req.CategoryID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrBudgetInvalidCategoryID
			}
			if _, err := s.db.GetCategoryByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrBudgetInvalidCategoryID
				}
				return nil, err
			}
		}
	}

	var period *ent.BudgetPeriod
	periodValue := ent.BudgetPeriod("")
	if req.Period != nil {
		parsed, ok := parseBudgetPeriod(strings.TrimSpace(*req.Period))
		if !ok {
			return nil, ErrBudgetPeriodInvalid
		}
		period = &parsed
		periodValue = parsed
	}

	startDate, err := parseBudgetDateString(req.StartDate)
	if err != nil {
		return nil, ErrBudgetDateInvalid
	}
	endDate, err := parseBudgetDateString(req.EndDate)
	if err != nil {
		return nil, ErrBudgetDateInvalid
	}

	if periodValue != "" || req.StartDate != nil || req.EndDate != nil {
		if periodValue == "" {
			current, err := s.db.GetBudgetByID(ctx, req.ID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrBudgetNotFound
				}
				return nil, err
			}
			periodValue = current.Period
		}

		resolvedStart, resolvedEnd, err := resolveBudgetDateRange(periodValue, startDate, endDate)
		if err != nil {
			return nil, err
		}
		startDate = resolvedStart
		endDate = resolvedEnd
	}

	item, err := s.db.UpdateBudget(ctx, req.ID, req.MemberID, req.CategoryID, req.Amount, period, startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, err
	}

	return &UpdateResponseService{ID: item.ID, MemberID: item.MemberID, CategoryID: item.CategoryID, Amount: item.Amount, SpentAmount: item.SpentAmount, Period: item.Period, StartDate: item.StartDate, EndDate: item.EndDate, CreatedAt: item.CreatedAt}, nil
}
