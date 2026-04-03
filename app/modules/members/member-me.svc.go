package members

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MeRequestService struct {
	MemberID string `json:"member_id"`
}

type MeResponseService struct {
	ID              uuid.UUID  `json:"id"`
	GenderID        *uuid.UUID `json:"gender_id"`
	PrefixID        *uuid.UUID `json:"prefix_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DisplayName     string     `json:"display_name"`
	Phone           string     `json:"phone"`
	ProfileImageURL string     `json:"profile_image_url"`
	Account         *MeAccount `json:"account"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLogin       *time.Time `json:"last_login"`
}

type MeAccount struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) InfoMeMember(ctx context.Context, req *MeRequestService) (*MeResponseService, error) {
	if _, err := uuid.Parse(req.MemberID); err != nil {
		return nil, ErrMemberUnauthorized
	}

	member, err := s.db.GetMemberByID(ctx, req.MemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	var account *MeAccount
	accounts, err := s.db.ListMemberAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range accounts {
		if item.MemberID != nil && *item.MemberID == member.ID {
			account = &MeAccount{ID: item.ID, Username: item.Username, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
			break
		}
	}

	return &MeResponseService{
		ID:              member.ID,
		GenderID:        member.GenderID,
		PrefixID:        member.PrefixID,
		FirstName:       member.FirstName,
		LastName:        member.LastName,
		DisplayName:     member.DisplayName,
		Phone:           member.Phone,
		ProfileImageURL: s.displayProfileImageURL(ctx, member.ProfileImageURL),
		Account:         account,
		CreatedAt:       member.CreatedAt,
		UpdatedAt:       member.UpdatedAt,
		LastLogin:       member.LastLogin,
	}, nil
}
