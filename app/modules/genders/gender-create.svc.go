package genders

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type CreateResponseService struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) CreateGender(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrGenderNameRequired
	}

	existing, err := s.db.ListGenders(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, item := range existing {
		if strings.EqualFold(strings.TrimSpace(item.Name), name) {
			return nil, ErrGenderAlreadyExists
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	gender, err := s.db.CreateGender(ctx, name, isActive)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{
		ID:        gender.ID,
		Name:      gender.Name,
		IsActive:  gender.IsActive,
		CreatedAt: gender.CreatedAt,
		UpdatedAt: gender.UpdatedAt,
	}, nil
}
