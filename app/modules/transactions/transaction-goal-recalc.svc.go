package transactions

import (
	"balance/app/modules/entities/ent"
	"context"
	"fmt"
	"math"
)

const goalRecalcEpsilon = 0.000001

func floatAlmostEqual(a float64, b float64) bool {
	return math.Abs(a-b) < goalRecalcEpsilon
}

func (s *Service) recalculateGoalsByWalletChanges(ctx context.Context, walletIDs []string, sourceID *string) error {
	memberIDs := make(map[string]struct{})

	for _, walletID := range walletIDs {
		if walletID == "" {
			continue
		}
		wallet, err := s.db.GetWalletByID(ctx, walletID)
		if err != nil {
			return err
		}
		if wallet.MemberID == nil {
			continue
		}
		memberIDs[wallet.MemberID.String()] = struct{}{}
	}

	changedWalletSet := make(map[string]struct{})
	for _, walletID := range walletIDs {
		if walletID == "" {
			continue
		}
		changedWalletSet[walletID] = struct{}{}
	}

	for memberID := range memberIDs {
		memberIDCopy := memberID
		goals, err := s.db.ListGoals(ctx, &memberIDCopy, nil, nil)
		if err != nil {
			return err
		}

		wallets, err := s.db.ListWallets(ctx, nil)
		if err != nil {
			return err
		}

		memberWalletBalances := make(map[string]float64)
		memberTotalWalletBalance := 0.0
		for _, wallet := range wallets {
			if wallet.MemberID == nil || wallet.MemberID.String() != memberID {
				continue
			}
			memberWalletBalances[wallet.ID.String()] = wallet.Balance
			memberTotalWalletBalance += wallet.Balance
		}

		for _, goal := range goals {
			if !goal.AutoTracking || goal.TrackingSourceType == nil {
				continue
			}

			nextAmount := goal.CurrentAmount
			switch *goal.TrackingSourceType {
			case ent.GoalTrackingSourceWallet:
				if goal.TrackingSourceID == nil {
					continue
				}
				trackedWalletID := goal.TrackingSourceID.String()
				if _, ok := changedWalletSet[trackedWalletID]; !ok {
					continue
				}
				nextAmount = memberWalletBalances[trackedWalletID]
			case ent.GoalTrackingSourceAllWallets:
				nextAmount = memberTotalWalletBalance
			default:
				continue
			}

			if floatAlmostEqual(goal.CurrentAmount, nextAmount) {
				continue
			}

			currentAmount := nextAmount
			_, err := s.db.UpdateGoal(ctx, goal.ID.String(), nil, nil, nil, &currentAmount, nil, nil, nil, nil, nil, nil, nil)
			if err != nil {
				return err
			}

			note := fmt.Sprintf("Auto recalculated from %s update", ent.GoalEntrySourceTypeTransaction)
			_, err = s.db.CreateGoalEntry(
				ctx,
				goal.ID.String(),
				&memberIDCopy,
				ent.GoalEntrySourceTypeTransaction,
				sourceID,
				goal.CurrentAmount,
				nextAmount,
				nextAmount-goal.CurrentAmount,
				note,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
