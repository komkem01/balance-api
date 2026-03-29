package categories

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	MemberID  *string `json:"member_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	IconName  string  `json:"icon_name"`
	ColorCode string  `json:"color_code"`
}

type CreateResponseService struct {
	ID        uuid.UUID        `json:"id"`
	MemberID  *uuid.UUID       `json:"member_id"`
	Name      string           `json:"name"`
	Type      ent.CategoryType `json:"type"`
	IconName  string           `json:"icon_name"`
	ColorCode string           `json:"color_code"`
	CreatedAt time.Time        `json:"created_at"`
}

func (s *Service) CreateCategory(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrCategoryNameRequired
	}
	categoryType, ok := parseCategoryType(strings.TrimSpace(req.Type))
	if !ok {
		return nil, ErrCategoryTypeInvalid
	}

	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrCategoryInvalidMemberID
			}
			if _, err := s.db.GetMemberByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrCategoryInvalidMemberID
				}
				return nil, err
			}
		}
	}

	item, err := s.db.CreateCategory(ctx, req.MemberID, strings.TrimSpace(req.Name), categoryType, strings.TrimSpace(req.IconName), strings.TrimSpace(req.ColorCode))
	if err != nil {
		return nil, err
	}
	return &CreateResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Type: item.Type, IconName: item.IconName, ColorCode: item.ColorCode, CreatedAt: item.CreatedAt}, nil
}
