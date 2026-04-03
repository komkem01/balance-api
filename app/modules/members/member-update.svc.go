package members

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID              string     `json:"id"`
	GenderID        *string    `json:"gender_id"`
	PrefixID        *string    `json:"prefix_id"`
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	DisplayName     *string    `json:"display_name"`
	Phone           *string    `json:"phone"`
	LastLogin       *time.Time `json:"last_login"`
	ProfileImageURL *string    `json:"profile_image_url"`
}

type UpdateResponseService struct {
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

func (s *Service) UpdateMember(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrMemberInvalidID
	}

	if req.GenderID == nil && req.PrefixID == nil && req.FirstName == nil && req.LastName == nil && req.DisplayName == nil && req.Phone == nil && req.LastLogin == nil && req.ProfileImageURL == nil {
		return nil, ErrMemberNoFieldsToUpdate
	}

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

	member, err := s.db.UpdateMember(
		ctx,
		req.ID,
		req.GenderID,
		req.PrefixID,
		req.FirstName,
		req.LastName,
		req.DisplayName,
		req.Phone,
		req.LastLogin,
		req.ProfileImageURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return &UpdateResponseService{
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
