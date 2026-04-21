package goals

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string
}

func (s *Service) DeleteGoal(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrGoalInvalidID
	}

	err := s.db.DeleteGoal(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGoalNotFound
		}
		return err
	}

	return nil
}
