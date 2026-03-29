package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ entitiesinf.CategoryEntity = (*Service)(nil)

func parseCategoryMemberID(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (s *Service) CreateCategory(ctx context.Context, memberID *string, name string, categoryType ent.CategoryType, iconName string, colorCode string) (*ent.CategoryEntity, error) {
	mid, err := parseCategoryMemberID(memberID)
	if err != nil {
		return nil, err
	}

	model := &ent.CategoryEntity{
		ID:        uuid.New(),
		MemberID:  mid,
		Name:      strings.TrimSpace(name),
		Type:      categoryType,
		IconName:  strings.TrimSpace(iconName),
		ColorCode: strings.TrimSpace(colorCode),
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id string) (*ent.CategoryEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.CategoryEntity{}
	if err := s.db.NewSelect().Model(model).Where("category.id = ?", uid).Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id string, memberID *string, name *string, categoryType *ent.CategoryType, iconName *string, colorCode *string) (*ent.CategoryEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.CategoryEntity{}
	if err := s.db.NewSelect().Model(model).Where("category.id = ?", uid).Scan(ctx); err != nil {
		return nil, err
	}

	if memberID != nil {
		mid, err := parseCategoryMemberID(memberID)
		if err != nil {
			return nil, err
		}
		model.MemberID = mid
	}
	if name != nil {
		model.Name = strings.TrimSpace(*name)
	}
	if categoryType != nil {
		model.Type = *categoryType
	}
	if iconName != nil {
		model.IconName = strings.TrimSpace(*iconName)
	}
	if colorCode != nil {
		model.ColorCode = strings.TrimSpace(*colorCode)
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("member_id", "name", "type", "icon_name", "color_code").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.db.NewDelete().Model(&ent.CategoryEntity{}).Where("id = ?", uid).Exec(ctx)
	return err
}

func (s *Service) ListCategories(ctx context.Context, memberID *string, categoryType *ent.CategoryType) ([]*ent.CategoryEntity, error) {
	items := make([]*ent.CategoryEntity, 0)
	q := s.db.NewSelect().Model(&items).Order("category.created_at DESC")
	if memberID != nil {
		mid, err := parseCategoryMemberID(memberID)
		if err != nil {
			return nil, err
		}
		if mid == nil {
			q = q.Where("category.member_id IS NULL")
		} else {
			q = q.Where("category.member_id = ?", *mid)
		}
	}
	if categoryType != nil {
		q = q.Where("category.type = ?", *categoryType)
	}
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
