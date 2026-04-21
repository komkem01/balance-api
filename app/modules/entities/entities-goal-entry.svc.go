package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ entitiesinf.GoalEntryEntity = (*Service)(nil)

func (s *Service) CreateGoalEntry(ctx context.Context, goalID string, memberID *string, sourceType ent.GoalEntrySourceType, sourceID *string, amountBefore float64, amountAfter float64, amountDelta float64, note string) (*ent.GoalEntryEntity, error) {
	gid, err := uuid.Parse(strings.TrimSpace(goalID))
	if err != nil {
		return nil, err
	}

	mid, err := parseGoalUUID(memberID)
	if err != nil {
		return nil, err
	}

	sid, err := parseGoalUUID(sourceID)
	if err != nil {
		return nil, err
	}

	model := &ent.GoalEntryEntity{
		ID:           uuid.New(),
		GoalID:       gid,
		MemberID:     mid,
		SourceType:   sourceType,
		SourceID:     sid,
		AmountBefore: amountBefore,
		AmountAfter:  amountAfter,
		AmountDelta:  amountDelta,
		Note:         strings.TrimSpace(note),
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) ListGoalEntries(ctx context.Context, goalID string, memberID *string) ([]*ent.GoalEntryEntity, error) {
	gid, err := uuid.Parse(strings.TrimSpace(goalID))
	if err != nil {
		return nil, err
	}

	items := make([]*ent.GoalEntryEntity, 0)
	q := s.db.NewSelect().Model(&items).Where("goal_entry.goal_id = ?", gid).Order("goal_entry.created_at DESC")

	if memberID != nil {
		mid, err := parseGoalUUID(memberID)
		if err != nil {
			return nil, err
		}
		if mid == nil {
			q = q.Where("goal_entry.member_id IS NULL")
		} else {
			q = q.Where("goal_entry.member_id = ?", *mid)
		}
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
