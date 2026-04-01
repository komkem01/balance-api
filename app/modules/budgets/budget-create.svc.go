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

type CreateRequestService struct {
	MemberID   *string `json:"member_id"`
	CategoryID *string `json:"category_id"`
	Amount     float64 `json:"amount"`
	Period     string  `json:"period"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
}

type CreateResponseService struct {
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

func parseBudgetDateString(value *string) (*time.Time, error) {
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

func (s *Service) CreateBudget(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	period, ok := parseBudgetPeriod(strings.TrimSpace(req.Period))
	if !ok {
		return nil, ErrBudgetPeriodInvalid
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

	startDate, err := parseBudgetDateString(req.StartDate)
	if err != nil {
		return nil, ErrBudgetDateInvalid
	}
	endDate, err := parseBudgetDateString(req.EndDate)
	if err != nil {
		return nil, ErrBudgetDateInvalid
	}
	startDate, endDate, err = resolveBudgetDateRange(period, startDate, endDate)
	if err != nil {
		return nil, err
	}

	item, err := s.db.CreateBudget(ctx, req.MemberID, req.CategoryID, req.Amount, period, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{ID: item.ID, MemberID: item.MemberID, CategoryID: item.CategoryID, Amount: item.Amount, SpentAmount: item.SpentAmount, Period: item.Period, StartDate: item.StartDate, EndDate: item.EndDate, CreatedAt: item.CreatedAt}, nil
}
