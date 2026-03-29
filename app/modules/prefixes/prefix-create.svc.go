package prefixes

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	GenderID string `json:"gender_id"`
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type CreateResponseService struct {
	ID        uuid.UUID `json:"id"`
	GenderID  uuid.UUID `json:"gender_id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) CreatePrefix(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	gid, err := uuid.Parse(req.GenderID)
	if err != nil {
		return nil, ErrPrefixInvalidGenderID
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrPrefixNameRequired
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	existing, err := s.db.ListPrefixes(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, item := range existing {
		if item.GenderID == gid && strings.EqualFold(strings.TrimSpace(item.Name), name) {
			return nil, ErrPrefixAlreadyExists
		}
	}

	prefix, err := s.db.CreatePrefix(ctx, req.GenderID, name, isActive)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{
		ID:        prefix.ID,
		GenderID:  prefix.GenderID,
		Name:      prefix.Name,
		IsActive:  prefix.IsActive,
		CreatedAt: prefix.CreatedAt,
		UpdatedAt: prefix.UpdatedAt,
	}, nil
}
