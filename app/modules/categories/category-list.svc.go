package categories

import (
	"context"
	"strings"
	"time"

	"balance/app/modules/entities/ent"
	"balance/app/utils/base"

	"github.com/google/uuid"
)

type ListRequestService struct {
	MemberID *string `json:"member_id"`
	Type     *string `json:"type"`
	Purpose  *string `json:"purpose"`
	Page     int     `json:"page"`
	Size     int     `json:"size"`
}

type ListItemService struct {
	ID        uuid.UUID            `json:"id"`
	MemberID  *uuid.UUID           `json:"member_id"`
	Name      string               `json:"name"`
	Type      ent.CategoryType     `json:"type"`
	Purpose   *ent.CategoryPurpose `json:"purpose"`
	IconName  string               `json:"icon_name"`
	ColorCode string               `json:"color_code"`
	CreatedAt time.Time            `json:"created_at"`
}

func (s *Service) ListCategory(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	var categoryType *ent.CategoryType
	if req.Type != nil {
		v := strings.TrimSpace(*req.Type)
		if v != "" {
			parsed, ok := parseCategoryType(v)
			if !ok {
				return nil, nil, ErrCategoryTypeInvalid
			}
			categoryType = &parsed
		}
	}

	var categoryPurpose *ent.CategoryPurpose
	if req.Purpose != nil {
		v := strings.TrimSpace(*req.Purpose)
		if v != "" {
			parsed, ok := parseCategoryPurpose(v)
			if !ok {
				return nil, nil, ErrCategoryPurposeInvalid
			}
			categoryPurpose = &parsed
		}
	}

	items, err := s.db.ListCategories(ctx, req.MemberID, categoryType, categoryPurpose)
	if err != nil {
		return nil, nil, err
	}
	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Type: item.Type, Purpose: item.Purpose, IconName: item.IconName, ColorCode: item.ColorCode, CreatedAt: item.CreatedAt})
	}

	page := int64(req.Page)
	if page < 1 {
		page = 1
	}
	size := int64(req.Size)
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	total := int64(len(res))
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return res[start:end], &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
