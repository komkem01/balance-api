package goals

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string
}

type InfoResponseService struct {
	ID                 uuid.UUID                   `json:"id"`
	MemberID           *uuid.UUID                  `json:"member_id"`
	Name               string                      `json:"name"`
	Type               ent.GoalType                `json:"type"`
	TargetAmount       float64                     `json:"target_amount"`
	StartAmount        float64                     `json:"start_amount"`
	CurrentAmount      float64                     `json:"current_amount"`
	StartDate          *time.Time                  `json:"start_date"`
	TargetDate         *time.Time                  `json:"target_date"`
	Status             ent.GoalStatus              `json:"status"`
	AutoTracking       bool                        `json:"auto_tracking"`
	TrackingSourceType *ent.GoalTrackingSourceType `json:"tracking_source_type"`
	TrackingSourceID   *uuid.UUID                  `json:"tracking_source_id"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

func (s *Service) InfoGoal(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrGoalInvalidID
	}

	item, err := s.db.GetGoalByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	res := &InfoResponseService{
		ID:                 item.ID,
		MemberID:           item.MemberID,
		Name:               item.Name,
		Type:               item.Type,
		TargetAmount:       item.TargetAmount,
		StartAmount:        item.StartAmount,
		CurrentAmount:      item.CurrentAmount,
		StartDate:          item.StartDate,
		TargetDate:         item.TargetDate,
		Status:             item.Status,
		AutoTracking:       item.AutoTracking,
		TrackingSourceType: item.TrackingSourceType,
		TrackingSourceID:   item.TrackingSourceID,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}

	resolvedCurrent, err := s.resolveAutoCurrentAmount(ctx, res)
	if err == nil {
		res.CurrentAmount = resolvedCurrent
	}

	return res, nil
}
