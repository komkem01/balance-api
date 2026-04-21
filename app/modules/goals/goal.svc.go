package goals

import (
	"balance/app/modules/entities/ent"
	"balance/internal/config"
	"context"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     GoalStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     GoalStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func isValidGoalType(goalType ent.GoalType) bool {
	switch goalType {
	case ent.GoalTypeSavings, ent.GoalTypeDebtPayoff:
		return true
	default:
		return false
	}
}

func isValidGoalStatus(status ent.GoalStatus) bool {
	switch status {
	case ent.GoalStatusActive, ent.GoalStatusCompleted, ent.GoalStatusPaused, ent.GoalStatusArchived:
		return true
	default:
		return false
	}
}

func isValidGoalSourceType(sourceType ent.GoalTrackingSourceType) bool {
	switch sourceType {
	case ent.GoalTrackingSourceWallet, ent.GoalTrackingSourceAllWallets, ent.GoalTrackingSourceLoan:
		return true
	default:
		return false
	}
}

func (s *Service) resolveAutoCurrentAmount(ctx context.Context, item *InfoResponseService) (float64, error) {
	if !item.AutoTracking || item.TrackingSourceType == nil {
		return item.CurrentAmount, nil
	}

	sourceType := *item.TrackingSourceType

	switch sourceType {
	case ent.GoalTrackingSourceWallet:
		if item.TrackingSourceID == nil {
			return item.CurrentAmount, nil
		}
		wallet, err := s.db.GetWalletByID(ctx, item.TrackingSourceID.String())
		if err != nil {
			return 0, err
		}
		if item.MemberID != nil && (wallet.MemberID == nil || wallet.MemberID.String() != item.MemberID.String()) {
			return 0, ErrGoalSourceMemberForbidden
		}
		return wallet.Balance, nil
	case ent.GoalTrackingSourceAllWallets:
		if item.MemberID == nil {
			return 0, nil
		}
		wallets, err := s.db.ListWallets(ctx, nil)
		if err != nil {
			return 0, err
		}
		total := 0.0
		for _, wallet := range wallets {
			if wallet.MemberID != nil && wallet.MemberID.String() == item.MemberID.String() {
				total += wallet.Balance
			}
		}
		return total, nil
	case ent.GoalTrackingSourceLoan:
		if item.TrackingSourceID == nil {
			return item.CurrentAmount, nil
		}
		loan, err := s.db.GetLoanByID(ctx, item.TrackingSourceID.String())
		if err != nil {
			return 0, err
		}
		if item.MemberID != nil && (loan.MemberID == nil || loan.MemberID.String() != item.MemberID.String()) {
			return 0, ErrGoalSourceMemberForbidden
		}

		paidAmount := loan.TotalAmount - loan.RemainingBalance
		if paidAmount < 0 {
			return 0, nil
		}
		return paidAmount, nil
	default:
		return item.CurrentAmount, nil
	}
}

func normalizeSourceID(value *string) *string {
	if value == nil {
		return nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil
	}
	return &v
}
