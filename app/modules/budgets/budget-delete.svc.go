package budgets

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteBudget(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrBudgetInvalidID
	}
	return s.db.DeleteBudget(ctx, req.ID)
}
