package memberaccounts

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	MemberID *string `json:"member_id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
}

type CreateResponseService struct {
	ID        uuid.UUID  `json:"id"`
	MemberID  *uuid.UUID `json:"member_id"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Service) CreateMemberAccount(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrMemberAccountInvalidMemberID
			}
		}
	}

	item, err := s.db.CreateMemberAccount(ctx, req.MemberID, strings.TrimSpace(req.Username), req.Password)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{
		ID:        item.ID,
		MemberID:  item.MemberID,
		Username:  item.Username,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
