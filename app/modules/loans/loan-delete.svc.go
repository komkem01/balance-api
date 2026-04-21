package loans

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string
}

func (s *Service) DeleteLoan(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrLoanInvalidID
	}
	err := s.db.DeleteLoan(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLoanNotFound
		}
		return err
	}
	return nil
}
