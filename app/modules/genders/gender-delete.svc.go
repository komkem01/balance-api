package genders

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteGender(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrGenderInvalidID
	}
	return s.db.DeleteGender(ctx, req.ID)
}
