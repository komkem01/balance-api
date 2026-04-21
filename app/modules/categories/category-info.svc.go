package categories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string `json:"id"`
}

type InfoResponseService struct {
	ID        uuid.UUID            `json:"id"`
	MemberID  *uuid.UUID           `json:"member_id"`
	Name      string               `json:"name"`
	Type      ent.CategoryType     `json:"type"`
	Purpose   *ent.CategoryPurpose `json:"purpose"`
	IconName  string               `json:"icon_name"`
	ColorCode string               `json:"color_code"`
	CreatedAt time.Time            `json:"created_at"`
}

func (s *Service) InfoCategory(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrCategoryInvalidID
	}
	item, err := s.db.GetCategoryByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &InfoResponseService{ID: item.ID, MemberID: item.MemberID, Name: item.Name, Type: item.Type, Purpose: item.Purpose, IconName: item.IconName, ColorCode: item.ColorCode, CreatedAt: item.CreatedAt}, nil
}
