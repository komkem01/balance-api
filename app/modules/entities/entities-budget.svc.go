package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.BudgetEntity = (*Service)(nil)

func parseBudgetID(value *string) (*uuid.UUID, error) {
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

func (s *Service) CreateBudget(ctx context.Context, memberID *string, categoryID *string, amount float64, period ent.BudgetPeriod, startDate *time.Time, endDate *time.Time) (*ent.BudgetEntity, error) {
	mid, err := parseBudgetID(memberID)
	if err != nil {
		return nil, err
	}
	cid, err := parseBudgetID(categoryID)
	if err != nil {
		return nil, err
	}

	model := &ent.BudgetEntity{
		ID:         uuid.New(),
		MemberID:   mid,
		CategoryID: cid,
		Amount:     amount,
		Period:     period,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) GetBudgetByID(ctx context.Context, id string) (*ent.BudgetEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.BudgetEntity{}
	if err := s.db.NewSelect().Model(model).Where("budget.id = ?", uid).Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateBudget(ctx context.Context, id string, memberID *string, categoryID *string, amount *float64, period *ent.BudgetPeriod, startDate *time.Time, endDate *time.Time) (*ent.BudgetEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.BudgetEntity{}
	if err := s.db.NewSelect().Model(model).Where("budget.id = ?", uid).Scan(ctx); err != nil {
		return nil, err
	}

	if memberID != nil {
		mid, err := parseBudgetID(memberID)
		if err != nil {
			return nil, err
		}
		model.MemberID = mid
	}
	if categoryID != nil {
		cid, err := parseBudgetID(categoryID)
		if err != nil {
			return nil, err
		}
		model.CategoryID = cid
	}
	if amount != nil {
		model.Amount = *amount
	}
	if period != nil {
		model.Period = *period
	}
	if startDate != nil {
		model.StartDate = startDate
	}
	if endDate != nil {
		model.EndDate = endDate
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("member_id", "category_id", "amount", "period", "start_date", "end_date").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) DeleteBudget(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.db.NewDelete().Model(&ent.BudgetEntity{}).Where("id = ?", uid).Exec(ctx)
	return err
}

func (s *Service) ListBudgets(ctx context.Context, memberID *string, categoryID *string, period *ent.BudgetPeriod) ([]*ent.BudgetEntity, error) {
	items := make([]*ent.BudgetEntity, 0)
	q := s.db.NewSelect().Model(&items).Order("budget.created_at DESC")

	if memberID != nil {
		mid, err := parseBudgetID(memberID)
		if err != nil {
			return nil, err
		}
		if mid == nil {
			q = q.Where("budget.member_id IS NULL")
		} else {
			q = q.Where("budget.member_id = ?", *mid)
		}
	}

	if categoryID != nil {
		cid, err := parseBudgetID(categoryID)
		if err != nil {
			return nil, err
		}
		if cid == nil {
			q = q.Where("budget.category_id IS NULL")
		} else {
			q = q.Where("budget.category_id = ?", *cid)
		}
	}

	if period != nil {
		q = q.Where("budget.period = ?", *period)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
