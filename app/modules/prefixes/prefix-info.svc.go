package prefixes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string `json:"id"`
}

type InfoResponseService struct {
	ID        uuid.UUID `json:"id"`
	GenderID  uuid.UUID `json:"gender_id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) InfoPrefix(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrPrefixInvalidID
	}

	prefix, err := s.db.GetPrefixByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPrefixNotFound
		}
		return nil, err
	}

	return &InfoResponseService{
		ID:        prefix.ID,
		GenderID:  prefix.GenderID,
		Name:      prefix.Name,
		IsActive:  prefix.IsActive,
		CreatedAt: prefix.CreatedAt,
		UpdatedAt: prefix.UpdatedAt,
	}, nil
}
