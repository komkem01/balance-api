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

type UpdateRequestService struct {
	ID        string  `json:"id"`
	MemberID  *string `json:"member_id"`
	Name      *string `json:"name"`
	Type      *string `json:"type"`
	IconName  *string `json:"icon_name"`
	ColorCode *string `json:"color_code"`
}

type UpdateResponseService struct {
	ID        uuid.UUID        `json:"id"`
	MemberID  *uuid.UUID       `json:"member_id"`
	Name      string           `json:"name"`
	Type      ent.CategoryType `json:"type"`
	IconName  string           `json:"icon_name"`
	ColorCode string           `json:"color_code"`
	CreatedAt time.Time        `json:"created_at"`
}

func (s *Service) UpdateCategory(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrCategoryInvalidID
	}
	if req.MemberID == nil && req.Name == nil && req.Type == nil && req.IconName == nil && req.ColorCode == nil {
		return nil, ErrCategoryNoFieldsToUpdate
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return nil, ErrCategoryNameRequired
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

	var categoryType *ent.CategoryType
	if req.Type != nil {
		parsed, ok := parseCategoryType(strings.TrimSpace(*req.Type))
		if !ok {
			return nil, ErrCategoryTypeInvalid
		}
		categoryType = &parsed
	}

	item, err := s.db.UpdateCategory(ctx, req.ID, req.MemberID, req.Name, categoryType, req.IconName, req.ColorCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &UpdateResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Type: item.Type, IconName: item.IconName, ColorCode: item.ColorCode, CreatedAt: item.CreatedAt}, nil
}
