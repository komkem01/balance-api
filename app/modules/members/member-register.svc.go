package members

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"balance/app/utils/hashing"

	"github.com/google/uuid"
)

type RegisterRequestService struct {
	GenderID    *string `json:"gender_id"`
	PrefixID    *string `json:"prefix_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Phone       string  `json:"phone"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
}

type RegisterResponseService struct {
	MemberID    uuid.UUID  `json:"member_id"`
	GenderID    *uuid.UUID `json:"gender_id"`
	PrefixID    *uuid.UUID `json:"prefix_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DisplayName string     `json:"display_name"`
	Phone       string     `json:"phone"`
	Username    string     `json:"username"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (s *Service) RegisterMember(ctx context.Context, req *RegisterRequestService) (*RegisterResponseService, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, ErrMemberUsernameRequired
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, ErrMemberPasswordRequired
	}

	accounts, err := s.db.ListMemberAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range accounts {
		if strings.EqualFold(strings.TrimSpace(item.Username), username) {
			return nil, ErrMemberUsernameExists
		}
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

	hashedPassword, err := hashing.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	member, err := s.db.CreateMemberWithAccount(
		ctx,
		req.GenderID,
		req.PrefixID,
		req.FirstName,
		req.LastName,
		req.DisplayName,
		req.Phone,
		username,
		string(hashedPassword),
	)
	if err != nil {
		return nil, err
	}

	return &RegisterResponseService{
		MemberID:    member.ID,
		GenderID:    member.GenderID,
		PrefixID:    member.PrefixID,
		FirstName:   member.FirstName,
		LastName:    member.LastName,
		DisplayName: member.DisplayName,
		Phone:       member.Phone,
		Username:    username,
		CreatedAt:   member.CreatedAt,
	}, nil
}
