package genders

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

type UpdateResponseService struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) UpdateGender(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrGenderInvalidID
	}
	if req.Name == nil && req.IsActive == nil {
		return nil, ErrGenderNoFieldsToUpdate
	}

	var name *string
	if req.Name != nil {
		value := strings.TrimSpace(*req.Name)
		if value == "" {
			return nil, ErrGenderNameRequired
		}

		existing, err := s.db.ListGenders(ctx, nil)
		if err != nil {
			return nil, err
		}
		for _, item := range existing {
			if item.ID.String() != req.ID && strings.EqualFold(strings.TrimSpace(item.Name), value) {
				return nil, ErrGenderAlreadyExists
			}
		}

		name = &value
	}

	gender, err := s.db.UpdateGender(ctx, req.ID, name, req.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGenderNotFound
		}
		return nil, err
	}

	return &UpdateResponseService{
		ID:        gender.ID,
		Name:      gender.Name,
		IsActive:  gender.IsActive,
		CreatedAt: gender.CreatedAt,
		UpdatedAt: gender.UpdatedAt,
	}, nil
}
