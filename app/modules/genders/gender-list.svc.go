package genders

import (
	"balance/app/utils/base"
	"context"
	"time"

	"github.com/google/uuid"
)

type ListRequestService struct {
	IsActive *bool `json:"is_active"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
}

type ListItemService struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) ListGender(ctx context.Context, req *ListRequestService) ([]*ListItemService, *base.ResponsePaginate, error) {
	items, err := s.db.ListGenders(ctx, req.IsActive)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*ListItemService, 0, len(items))
	for _, item := range items {
		res = append(res, &ListItemService{
			ID:        item.ID,
			Name:      item.Name,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
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
