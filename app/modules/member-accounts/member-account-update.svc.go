package memberaccounts

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
	MemberID *string `json:"member_id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UpdateResponseService struct {
	ID        uuid.UUID  `json:"id"`
	MemberID  *uuid.UUID `json:"member_id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Service) UpdateMemberAccount(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrMemberAccountInvalidID
	}
	if req.MemberID == nil && req.Username == nil && req.Password == nil {
		return nil, ErrMemberAccountNoFieldsToUpdate
	}

	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrMemberAccountInvalidMemberID
			}
		}
	}

	item, err := s.db.UpdateMemberAccount(ctx, req.ID, req.MemberID, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberAccountNotFound
		}
		return nil, err
	}

	return &UpdateResponseService{
		ID:        item.ID,
		MemberID:  item.MemberID,
		Username:  item.Username,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
