package members

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	GenderID    *string `json:"gender_id"`
	PrefixID    *string `json:"prefix_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Phone       string  `json:"phone"`
}

type CreateResponseService struct {
	ID              uuid.UUID  `json:"id"`
	GenderID        *uuid.UUID `json:"gender_id"`
	PrefixID        *uuid.UUID `json:"prefix_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DisplayName     string     `json:"display_name"`
	Phone           string     `json:"phone"`
	ProfileImageURL string     `json:"profile_image_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLogin       *time.Time `json:"last_login"`
}

func (s *Service) CreateMember(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if req.GenderID != nil {
		v := strings.TrimSpace(*req.GenderID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrMemberInvalidGenderID
			}
			if _, err := s.db.GetGenderByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrMemberInvalidGenderID
				}
				return nil, err
			}
		}
	}

	if req.PrefixID != nil {
		v := strings.TrimSpace(*req.PrefixID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrMemberInvalidPrefixID
			}
			if _, err := s.db.GetPrefixByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrMemberInvalidPrefixID
				}
				return nil, err
			}
		}
	}

	member, err := s.db.CreateMember(
		ctx,
		req.GenderID,
		req.PrefixID,
		strings.TrimSpace(req.FirstName),
		strings.TrimSpace(req.LastName),
		strings.TrimSpace(req.DisplayName),
		strings.TrimSpace(req.Phone),
	)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{
		ID:              member.ID,
		GenderID:        member.GenderID,
		PrefixID:        member.PrefixID,
		FirstName:       member.FirstName,
		LastName:        member.LastName,
		DisplayName:     member.DisplayName,
		Phone:           member.Phone,
		ProfileImageURL: s.displayProfileImageURL(ctx, member.ProfileImageURL),
		CreatedAt:       member.CreatedAt,
		UpdatedAt:       member.UpdatedAt,
		LastLogin:       member.LastLogin,
	}, nil
}
