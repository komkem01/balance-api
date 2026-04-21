package goals

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	MemberID           *string                     `json:"member_id"`
	Name               string                      `json:"name"`
	Type               ent.GoalType                `json:"type"`
	TargetAmount       float64                     `json:"target_amount"`
	StartAmount        float64                     `json:"start_amount"`
	CurrentAmount      float64                     `json:"current_amount"`
	StartDate          *string                     `json:"start_date"`
	TargetDate         *string                     `json:"target_date"`
	Status             ent.GoalStatus              `json:"status"`
	AutoTracking       *bool                       `json:"auto_tracking"`
	TrackingSourceType *ent.GoalTrackingSourceType `json:"tracking_source_type"`
	TrackingSourceID   *string                     `json:"tracking_source_id"`
	DepositWalletID    *string                     `json:"deposit_wallet_id"`
}

type CreateResponseService struct {
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
	DepositWalletID    *uuid.UUID                  `json:"deposit_wallet_id"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

func parseGoalDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Service) resolveSourceCurrentAmount(ctx context.Context, memberID *string, goalType ent.GoalType, sourceType *ent.GoalTrackingSourceType, sourceID *string) (float64, error) {
	if sourceType == nil {
		return 0, nil
	}

	if !isValidGoalSourceType(*sourceType) {
		return 0, ErrGoalSourceTypeInvalid
	}

	sid := normalizeSourceID(sourceID)

	switch *sourceType {
	case ent.GoalTrackingSourceWallet:
		if sid == nil {
			return 0, ErrGoalSourceIDRequired
		}
		wallet, err := s.db.GetWalletByID(ctx, *sid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, ErrGoalSourceIDInvalid
			}
			return 0, err
		}
		if memberID != nil {
			mid := strings.TrimSpace(*memberID)
			if wallet.MemberID == nil || wallet.MemberID.String() != mid {
				return 0, ErrGoalSourceMemberForbidden
			}
		}
		return wallet.Balance, nil
	case ent.GoalTrackingSourceAllWallets:
		if memberID == nil {
			return 0, nil
		}
		mid := strings.TrimSpace(*memberID)
		wallets, err := s.db.ListWallets(ctx, nil)
		if err != nil {
			return 0, err
		}
		total := 0.0
		for _, wallet := range wallets {
			if wallet.MemberID != nil && wallet.MemberID.String() == mid {
				total += wallet.Balance
			}
		}
		return total, nil
	case ent.GoalTrackingSourceLoan:
		if sid == nil {
			return 0, ErrGoalSourceIDRequired
		}
		loan, err := s.db.GetLoanByID(ctx, *sid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, ErrGoalSourceIDInvalid
			}
			return 0, err
		}
		if memberID != nil {
			mid := strings.TrimSpace(*memberID)
			if loan.MemberID == nil || loan.MemberID.String() != mid {
				return 0, ErrGoalSourceMemberForbidden
			}
		}
		if goalType == ent.GoalTypeDebtPayoff {
			paid := loan.TotalAmount - loan.RemainingBalance
			if paid < 0 {
				return 0, nil
			}
			return paid, nil
		}
		return loan.TotalAmount - loan.RemainingBalance, nil
	default:
		return 0, nil
	}
}

func (s *Service) CreateGoal(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrGoalNameRequired
	}
	if !isValidGoalType(req.Type) {
		return nil, ErrGoalTypeInvalid
	}
	if req.TargetAmount < 0 {
		return nil, ErrGoalTargetAmountInvalid
	}

	status := req.Status
	if status == "" {
		status = ent.GoalStatusActive
	}
	if !isValidGoalStatus(status) {
		return nil, ErrGoalStatusInvalid
	}

	autoTracking := true
	if req.AutoTracking != nil {
		autoTracking = *req.AutoTracking
	}

	if normalizedSourceID := normalizeSourceID(req.TrackingSourceID); normalizedSourceID != nil {
		if _, err := uuid.Parse(*normalizedSourceID); err != nil {
			return nil, ErrGoalSourceIDInvalid
		}
	}

	if req.MemberID != nil {
		memberID := strings.TrimSpace(*req.MemberID)
		if memberID != "" {
			if _, err := uuid.Parse(memberID); err != nil {
				return nil, ErrGoalInvalidMemberID
			}
			if _, err := s.db.GetMemberByID(ctx, memberID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrGoalInvalidMemberID
				}
				return nil, err
			}
		}
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

		if req.MemberID != nil {
			mid := strings.TrimSpace(*req.MemberID)
			if wallet.MemberID == nil || wallet.MemberID.String() != mid {
				return nil, ErrGoalDepositWalletForbidden
			}
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

	currentAmount := req.CurrentAmount
	if autoTracking {
		resolvedCurrent, err := s.resolveSourceCurrentAmount(ctx, req.MemberID, req.Type, req.TrackingSourceType, req.TrackingSourceID)
		if err != nil {
			return nil, err
		}
		currentAmount = resolvedCurrent
	}

	normalizedSourceID := normalizeSourceID(req.TrackingSourceID)
	normalizedDepositWalletID := normalizeSourceID(req.DepositWalletID)

	item, err := s.db.CreateGoal(ctx, req.MemberID, req.Name, req.Type, req.TargetAmount, req.StartAmount, currentAmount, startDate, targetDate, status, autoTracking, req.TrackingSourceType, normalizedSourceID, normalizedDepositWalletID)
	if err != nil {
		return nil, err
	}

	return &CreateResponseService{
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
		DepositWalletID:    item.DepositWalletID,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}, nil
}
