package goals

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID                 string                      `json:"id"`
	Name               *string                     `json:"name"`
	TargetAmount       *float64                    `json:"target_amount"`
	StartAmount        *float64                    `json:"start_amount"`
	CurrentAmount      *float64                    `json:"current_amount"`
	StartDate          *string                     `json:"start_date"`
	TargetDate         *string                     `json:"target_date"`
	Status             *ent.GoalStatus             `json:"status"`
	AutoTracking       *bool                       `json:"auto_tracking"`
	TrackingSourceType *ent.GoalTrackingSourceType `json:"tracking_source_type"`
	TrackingSourceID   *string                     `json:"tracking_source_id"`
	DepositWalletID    *string                     `json:"deposit_wallet_id"`
}

func (s *Service) UpdateGoal(ctx context.Context, req *UpdateRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrGoalInvalidID
	}

	if req.Name == nil && req.TargetAmount == nil && req.StartAmount == nil && req.CurrentAmount == nil && req.StartDate == nil && req.TargetDate == nil && req.Status == nil && req.AutoTracking == nil && req.TrackingSourceType == nil && req.TrackingSourceID == nil && req.DepositWalletID == nil {
		return nil, ErrGoalNoFieldsToUpdate
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return nil, ErrGoalNameRequired
	}

	if req.TargetAmount != nil && *req.TargetAmount < 0 {
		return nil, ErrGoalTargetAmountInvalid
	}

	if req.Status != nil && !isValidGoalStatus(*req.Status) {
		return nil, ErrGoalStatusInvalid
	}

	if req.TrackingSourceType != nil && !isValidGoalSourceType(*req.TrackingSourceType) {
		return nil, ErrGoalSourceTypeInvalid
	}

	if normalizedSourceID := normalizeSourceID(req.TrackingSourceID); normalizedSourceID != nil {
		if _, err := uuid.Parse(*normalizedSourceID); err != nil {
			return nil, ErrGoalSourceIDInvalid
		}
	}

	startDate, err := parseGoalDate(req.StartDate)
	if err != nil {
		return nil, ErrGoalStartDateInvalid
	}
	targetDate, err := parseGoalDate(req.TargetDate)
	if err != nil {
		return nil, ErrGoalTargetDateInvalid
	}

	item, err := s.db.GetGoalByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	if normalizedDepositWalletID := normalizeSourceID(req.DepositWalletID); normalizedDepositWalletID != nil {
		if _, err := uuid.Parse(*normalizedDepositWalletID); err != nil {
			return nil, ErrGoalDepositWalletInvalid
		}

		wallet, err := s.db.GetWalletByID(ctx, *normalizedDepositWalletID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrGoalDepositWalletInvalid
			}
			return nil, err
		}

		if item.MemberID != nil {
			mid := item.MemberID.String()
			if wallet.MemberID == nil || wallet.MemberID.String() != mid {
				return nil, ErrGoalDepositWalletForbidden
			}
		}
	}

	normalizedSourceID := normalizeSourceID(req.TrackingSourceID)

	autoTracking := item.AutoTracking
	if req.AutoTracking != nil {
		autoTracking = *req.AutoTracking
	}

	if autoTracking {
		effectiveSourceType := item.TrackingSourceType
		if req.TrackingSourceType != nil {
			effectiveSourceType = req.TrackingSourceType
		}

		effectiveSourceID := req.TrackingSourceID
		if effectiveSourceID == nil && item.TrackingSourceID != nil {
			t := item.TrackingSourceID.String()
			effectiveSourceID = &t
		}

		memberID := ""
		if item.MemberID != nil {
			memberID = item.MemberID.String()
		}

		currentAmount, err := s.resolveSourceCurrentAmount(ctx, &memberID, item.Type, effectiveSourceType, effectiveSourceID)
		if err != nil {
			return nil, err
		}
		req.CurrentAmount = &currentAmount
	}

	normalizedDepositWalletID := normalizeSourceID(req.DepositWalletID)

	updated, err := s.db.UpdateGoal(ctx, req.ID, req.Name, req.TargetAmount, req.StartAmount, req.CurrentAmount, startDate, targetDate, req.Status, req.AutoTracking, req.TrackingSourceType, normalizedSourceID, normalizedDepositWalletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, err
	}

	res := &InfoResponseService{
		ID:                 updated.ID,
		MemberID:           updated.MemberID,
		Name:               updated.Name,
		Type:               updated.Type,
		TargetAmount:       updated.TargetAmount,
		StartAmount:        updated.StartAmount,
		CurrentAmount:      updated.CurrentAmount,
		StartDate:          updated.StartDate,
		TargetDate:         updated.TargetDate,
		Status:             updated.Status,
		AutoTracking:       updated.AutoTracking,
		TrackingSourceType: updated.TrackingSourceType,
		TrackingSourceID:   updated.TrackingSourceID,
		DepositWalletID:    updated.DepositWalletID,
		CreatedAt:          updated.CreatedAt,
		UpdatedAt:          updated.UpdatedAt,
	}

	return res, nil
}
