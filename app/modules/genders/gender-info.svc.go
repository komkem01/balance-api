package genders

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
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) InfoGender(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrGenderInvalidID
	}

	gender, err := s.db.GetGenderByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGenderNotFound
		}
		return nil, err
	}

	return &InfoResponseService{
		ID:        gender.ID,
		Name:      gender.Name,
		IsActive:  gender.IsActive,
		CreatedAt: gender.CreatedAt,
		UpdatedAt: gender.UpdatedAt,
	}, nil
}
