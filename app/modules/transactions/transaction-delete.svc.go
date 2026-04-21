package transactions

import (
	"balance/app/modules/entities"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteTransaction(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrTransactionInvalidID
	}

	item, err := s.db.GetTransactionByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransactionNotFound
		}
		return err
	}
	err = s.db.DeleteTransactionWithWalletAdjust(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransactionNotFound
		}
		if errors.Is(err, entities.ErrWalletBalanceInsufficient) {
			return ErrTransactionInsufficientFunds
		}
		return err
	}

	if item.WalletID != nil {
		sourceID := item.ID.String()
		if err := s.recalculateGoalsByWalletChanges(ctx, []string{item.WalletID.String()}, &sourceID); err != nil {
			return err
		}
	}

	return nil
}
