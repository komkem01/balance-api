package budgets

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string `json:"id"`
}

type InfoResponseService struct {
	ID         uuid.UUID        `json:"id"`
	MemberID   *uuid.UUID       `json:"member_id"`
	CategoryID *uuid.UUID       `json:"category_id"`
	Amount     float64          `json:"amount"`
	Period     ent.BudgetPeriod `json:"period"`
	StartDate  *time.Time       `json:"start_date"`
	EndDate    *time.Time       `json:"end_date"`
	CreatedAt  time.Time        `json:"created_at"`
}

func (s *Service) InfoBudget(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrBudgetInvalidID
	}
	item, err := s.db.GetBudgetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, err
	}

	return &InfoResponseService{ID: item.ID, MemberID: item.MemberID, CategoryID: item.CategoryID, Amount: item.Amount, Period: item.Period, StartDate: item.StartDate, EndDate: item.EndDate, CreatedAt: item.CreatedAt}, nil
}
