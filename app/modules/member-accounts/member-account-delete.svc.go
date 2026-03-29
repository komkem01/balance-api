package memberaccounts

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteMemberAccount(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrMemberAccountInvalidID
	}
	return s.db.DeleteMemberAccount(ctx, req.ID)
}
