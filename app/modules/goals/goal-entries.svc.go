package goals

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type ListEntriesRequestService struct {
	GoalID   string
	MemberID *string
}

type EntryResponseService struct {
	ID           uuid.UUID               `json:"id"`
	GoalID       uuid.UUID               `json:"goal_id"`
	MemberID     *uuid.UUID              `json:"member_id"`
	SourceType   ent.GoalEntrySourceType `json:"source_type"`
	SourceID     *uuid.UUID              `json:"source_id"`
	AmountBefore float64                 `json:"amount_before"`
	AmountAfter  float64                 `json:"amount_after"`
	AmountDelta  float64                 `json:"amount_delta"`
	Note         string                  `json:"note"`
	CreatedAt    time.Time               `json:"created_at"`
}

func (s *Service) ListGoalEntries(ctx context.Context, req *ListEntriesRequestService) ([]*EntryResponseService, error) {
	if _, err := uuid.Parse(req.GoalID); err != nil {
		return nil, ErrGoalInvalidID
	}

	if _, err := s.db.GetGoalByID(ctx, req.GoalID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	items, err := s.db.ListGoalEntries(ctx, req.GoalID, req.MemberID)
	if err != nil {
		return nil, err
	}

	res := make([]*EntryResponseService, 0, len(items))
	for _, item := range items {
		res = append(res, &EntryResponseService{
			ID:           item.ID,
			GoalID:       item.GoalID,
			MemberID:     item.MemberID,
			SourceType:   item.SourceType,
			SourceID:     item.SourceID,
			AmountBefore: item.AmountBefore,
			AmountAfter:  item.AmountAfter,
			AmountDelta:  item.AmountDelta,
			Note:         item.Note,
			CreatedAt:    item.CreatedAt,
		})
	}

	return res, nil
}
