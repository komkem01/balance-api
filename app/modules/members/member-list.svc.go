package members

import (
	"balance/app/utils/base"
	"context"
	"time"

	"github.com/google/uuid"
)

type ListRequestService struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListItemService struct {
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

func (s *Service) ListMember(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	items, err := s.db.ListMembers(ctx)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{
			ID:          item.ID,
			GenderID:    item.GenderID,
			PrefixID:    item.PrefixID,
			FirstName:   item.FirstName,
			LastName:    item.LastName,
			DisplayName: item.DisplayName,
			Phone:       item.Phone,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			LastLogin:   item.LastLogin,
		})
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

	paged := res[start:end]
	return paged, &base.ResponsePaginate{Page: page, Size: size, Total: total}, nil
}
