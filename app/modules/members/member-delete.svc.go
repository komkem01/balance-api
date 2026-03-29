package members

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteMember(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrMemberInvalidID
	}
	return s.db.DeleteMember(ctx, req.ID)
}
