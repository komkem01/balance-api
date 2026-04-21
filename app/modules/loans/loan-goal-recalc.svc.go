package loans

import (
	"balance/app/modules/entities/ent"
	"context"
	"fmt"
	"math"
)

const loanGoalRecalcEpsilon = 0.000001

func loanFloatAlmostEqual(a float64, b float64) bool {
	return math.Abs(a-b) < loanGoalRecalcEpsilon
}

func (s *Service) recalculateGoalsByLoanUpdate(ctx context.Context, item *UpdateResponseService) error {
	if item == nil || item.MemberID == nil {
		return nil
	}

	memberID := item.MemberID.String()
	goals, err := s.db.ListGoals(ctx, &memberID, nil, nil)
	if err != nil {
		return err
	}

	loanID := item.ID.String()
	nextAmount := item.TotalAmount - item.RemainingBalance
	if nextAmount < 0 {
		nextAmount = 0
	}

	for _, goal := range goals {
		if !goal.AutoTracking || goal.TrackingSourceType == nil || *goal.TrackingSourceType != ent.GoalTrackingSourceLoan {
			continue
		}
		if goal.TrackingSourceID == nil || goal.TrackingSourceID.String() != loanID {
			continue
		}
		if loanFloatAlmostEqual(goal.CurrentAmount, nextAmount) {
			continue
		}

		currentAmount := nextAmount
		_, err := s.db.UpdateGoal(ctx, goal.ID.String(), nil, nil, nil, &currentAmount, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		sourceID := loanID
		note := fmt.Sprintf("GOAL_ENTRY|source=loan|action=auto_recalculate|loan_id=%s", loanID)
		_, err = s.db.CreateGoalEntry(
			ctx,
			goal.ID.String(),
			&memberID,
			ent.GoalEntrySourceTypeLoan,
			&sourceID,
			goal.CurrentAmount,
			nextAmount,
			nextAmount-goal.CurrentAmount,
			note,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
