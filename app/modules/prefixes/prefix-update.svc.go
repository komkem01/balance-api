package prefixes

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

func (s *Service) UpdatePrefix(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrPrefixInvalidID
	}
	if req.Name == nil && req.IsActive == nil {
		return nil, ErrPrefixNoFieldsToUpdate
	}

	current, err := s.db.GetPrefixByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPrefixNotFound
		}
		return nil, err
	}

	var name *string
	if req.Name != nil {
		value := strings.TrimSpace(*req.Name)
		if value == "" {
			return nil, ErrPrefixNameRequired
		}

		existing, err := s.db.ListPrefixes(ctx, nil)
		if err != nil {
			return nil, err
		}
		for _, item := range existing {
			if item.ID != current.ID && item.GenderID == current.GenderID && strings.EqualFold(strings.TrimSpace(item.Name), value) {
				return nil, ErrPrefixAlreadyExists
			}
		}

		name = &value
	}

	prefix, err := s.db.UpdatePrefix(ctx, req.ID, name, req.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPrefixNotFound
		}
		return nil, err
	}

	return &UpdateResponseService{
		ID:        prefix.ID,
		Name:      prefix.Name,
		IsActive:  prefix.IsActive,
		CreatedAt: prefix.CreatedAt,
		UpdatedAt: prefix.UpdatedAt,
	}, nil
}
