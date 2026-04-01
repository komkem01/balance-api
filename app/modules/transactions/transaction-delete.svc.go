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
	err := s.db.DeleteTransactionWithWalletAdjust(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransactionNotFound
		}
		if errors.Is(err, entities.ErrWalletBalanceInsufficient) {
			return ErrTransactionInsufficientFunds
		}
	}
	return err
}
