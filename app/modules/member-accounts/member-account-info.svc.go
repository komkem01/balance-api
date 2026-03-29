package memberaccounts

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
	ID        uuid.UUID  `json:"id"`
	MemberID  *uuid.UUID `json:"member_id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Service) InfoMemberAccount(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrMemberAccountInvalidID
	}

	item, err := s.db.GetMemberAccountByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberAccountNotFound
		}
		return nil, err
	}

	return &InfoResponseService{
		ID:        item.ID,
		MemberID:  item.MemberID,
		Username:  item.Username,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
