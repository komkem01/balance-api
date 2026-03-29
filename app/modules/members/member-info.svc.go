package members

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
	ID          uuid.UUID  `json:"id"`
	GenderID    *uuid.UUID `json:"gender_id"`
	PrefixID    *uuid.UUID `json:"prefix_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DisplayName string     `json:"display_name"`
	Phone       string     `json:"phone"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLogin   *time.Time `json:"last_login"`
}

func (s *Service) InfoMember(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrMemberInvalidID
	}

	member, err := s.db.GetMemberByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return &InfoResponseService{
		ID:          member.ID,
		GenderID:    member.GenderID,
		PrefixID:    member.PrefixID,
		FirstName:   member.FirstName,
		LastName:    member.LastName,
		DisplayName: member.DisplayName,
		Phone:       member.Phone,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
		LastLogin:   member.LastLogin,
	}, nil
}
