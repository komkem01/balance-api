package members

import (
	"context"
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

	member, err := s.CreateMember(ctx, &CreateRequestService{
		GenderID:    req.GenderID,
		PrefixID:    req.PrefixID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Phone:       req.Phone,
	})
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashing.HashPassword(req.Password)
	if err != nil {
		_ = s.db.DeleteMember(ctx, member.ID.String())
		return nil, err
	}

	memberID := member.ID.String()
	if _, err := s.db.CreateMemberAccount(ctx, &memberID, username, string(hashedPassword)); err != nil {
		_ = s.db.DeleteMember(ctx, member.ID.String())
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
