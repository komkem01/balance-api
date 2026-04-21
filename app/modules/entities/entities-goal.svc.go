package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.GoalEntity = (*Service)(nil)

func parseGoalUUID(value *string) (*uuid.UUID, error) {
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

func (s *Service) CreateGoal(ctx context.Context, memberID *string, name string, goalType ent.GoalType, targetAmount float64, startAmount float64, currentAmount float64, startDate *time.Time, targetDate *time.Time, status ent.GoalStatus, autoTracking bool, trackingSourceType *ent.GoalTrackingSourceType, trackingSourceID *string) (*ent.GoalEntity, error) {
	mid, err := parseGoalUUID(memberID)
	if err != nil {
		return nil, err
	}
	sid, err := parseGoalUUID(trackingSourceID)
	if err != nil {
		return nil, err
	}

	model := &ent.GoalEntity{
		ID:                 uuid.New(),
		MemberID:           mid,
		Name:               strings.TrimSpace(name),
		Type:               goalType,
		TargetAmount:       targetAmount,
		StartAmount:        startAmount,
		CurrentAmount:      currentAmount,
		StartDate:          startDate,
		TargetDate:         targetDate,
		Status:             status,
		AutoTracking:       autoTracking,
		TrackingSourceType: trackingSourceType,
		TrackingSourceID:   sid,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) GetGoalByID(ctx context.Context, id string) (*ent.GoalEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.GoalEntity{}
	if err := s.db.NewSelect().Model(model).Where("goal.id = ?", uid).Where("goal.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateGoal(ctx context.Context, id string, name *string, targetAmount *float64, startAmount *float64, currentAmount *float64, startDate *time.Time, targetDate *time.Time, status *ent.GoalStatus, autoTracking *bool, trackingSourceType *ent.GoalTrackingSourceType, trackingSourceID *string) (*ent.GoalEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.GoalEntity{}
	if err := s.db.NewSelect().Model(model).Where("goal.id = ?", uid).Where("goal.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}

	if name != nil {
		model.Name = strings.TrimSpace(*name)
	}
	if targetAmount != nil {
		model.TargetAmount = *targetAmount
	}
	if startAmount != nil {
		model.StartAmount = *startAmount
	}
	if currentAmount != nil {
		model.CurrentAmount = *currentAmount
	}
	if startDate != nil {
		model.StartDate = startDate
	}
	if targetDate != nil {
		model.TargetDate = targetDate
	}
	if status != nil {
		model.Status = *status
	}
	if autoTracking != nil {
		model.AutoTracking = *autoTracking
	}
	if trackingSourceType != nil {
		model.TrackingSourceType = trackingSourceType
	}
	if trackingSourceID != nil {
		sid, err := parseGoalUUID(trackingSourceID)
		if err != nil {
			return nil, err
		}
		model.TrackingSourceID = sid
	}

	model.UpdatedAt = time.Now()

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("name", "target_amount", "start_amount", "current_amount", "start_date", "target_date", "status", "auto_tracking", "tracking_source_type", "tracking_source_id", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) DeleteGoal(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	model := &ent.GoalEntity{ID: uid}
	_, err = s.db.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (s *Service) ListGoals(ctx context.Context, memberID *string, status *ent.GoalStatus, goalType *ent.GoalType) ([]*ent.GoalEntity, error) {
	items := make([]*ent.GoalEntity, 0)
	q := s.db.NewSelect().Model(&items).Where("goal.deleted_at IS NULL").Order("goal.created_at DESC")

	if memberID != nil {
		mid, err := parseGoalUUID(memberID)
		if err != nil {
			return nil, err
		}
		if mid != nil {
			q = q.Where("goal.member_id = ?", *mid)
		}
	}

	if status != nil {
		q = q.Where("goal.status = ?", *status)
	}

	if goalType != nil {
		q = q.Where("goal.type = ?", *goalType)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
