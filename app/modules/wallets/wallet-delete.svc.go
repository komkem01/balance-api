package wallets

import (
	"context"

	"github.com/google/uuid"
)

type DeleteRequestService struct {
	ID string `json:"id"`
}

func (s *Service) DeleteWallet(ctx context.Context, req *DeleteRequestService) error {
	if _, err := uuid.Parse(req.ID); err != nil {
		return ErrWalletInvalidID
	}
	return s.db.DeleteWallet(ctx, req.ID)
}
