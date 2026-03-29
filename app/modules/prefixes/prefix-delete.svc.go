package prefixes

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeletePrefix(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrPrefixInvalidID
	}
	return s.db.DeletePrefix(ctx, req.ID)
}
