package transactions

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteTransaction(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrTransactionInvalidID
	}
	return s.db.DeleteTransaction(ctx, req.ID)
}
