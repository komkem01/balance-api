package transactions

import (
	"balance/app/modules/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransferRequestService struct {
	FromWalletID    string  `json:"from_wallet_id"`
	ToWalletID      string  `json:"to_wallet_id"`
	CategoryID      *string `json:"category_id"`
	Amount          float64 `json:"amount"`
	TransactionDate *string `json:"transaction_date"`
	Note            string  `json:"note"`
}

type TransferResponseService struct {
	FromTransaction *CreateResponseService `json:"from_transaction"`
	ToTransaction   *CreateResponseService `json:"to_transaction"`
}

const transferNotePrefix = "__transfer__"

func buildTransferNote(ref string, direction string, counterpartyWalletID string, userNote string) string {
	cleanNote := strings.ReplaceAll(strings.TrimSpace(userNote), "|", "/")
	return fmt.Sprintf("%s|%s|%s|%s|%s", transferNotePrefix, ref, direction, strings.TrimSpace(counterpartyWalletID), cleanNote)
}

func parseTransferDateString(value *string) (*time.Time, error) {
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

func (s *Service) TransferBetweenWallets(ctx context.Context, req *TransferRequestService) (*TransferResponseService, error) {
	if req.Amount <= 0 {
		return nil, ErrTransactionAmountInvalid
	}

	fromWalletID := strings.TrimSpace(req.FromWalletID)
	toWalletID := strings.TrimSpace(req.ToWalletID)

	fromUUID, err := uuid.Parse(fromWalletID)
	if err != nil {
		return nil, ErrTransactionTransferInvalidFromWalletID
	}
	toUUID, err := uuid.Parse(toWalletID)
	if err != nil {
		return nil, ErrTransactionTransferInvalidToWalletID
	}
	if fromUUID == toUUID {
		return nil, ErrTransactionTransferSameWallet
	}

	if _, err := s.db.GetWalletByID(ctx, fromWalletID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionTransferInvalidFromWalletID
		}
		return nil, err
	}
	if _, err := s.db.GetWalletByID(ctx, toWalletID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionTransferInvalidToWalletID
		}
		return nil, err
	}

	if req.CategoryID != nil {
		categoryID := strings.TrimSpace(*req.CategoryID)
		if categoryID != "" {
			if _, err := uuid.Parse(categoryID); err != nil {
				return nil, ErrTransactionInvalidCategoryID
			}
			if _, err := s.db.GetCategoryByID(ctx, categoryID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrTransactionInvalidCategoryID
				}
				return nil, err
			}
		}
	}

	transactionDate, err := parseTransferDateString(req.TransactionDate)
	if err != nil {
		return nil, ErrTransactionDateInvalid
	}

	transferRef := uuid.NewString()
	fromNote := buildTransferNote(transferRef, "out", toWalletID, req.Note)
	toNote := buildTransferNote(transferRef, "in", fromWalletID, req.Note)

	fromTx, toTx, err := s.db.CreateTransferTransactionsWithWalletAdjust(
		ctx,
		fromWalletID,
		toWalletID,
		req.CategoryID,
		req.Amount,
		transactionDate,
		fromNote,
		toNote,
	)
	if err != nil {
		if errors.Is(err, entities.ErrWalletBalanceInsufficient) {
			return nil, ErrTransactionInsufficientFunds
		}
		return nil, err
	}

	return &TransferResponseService{
		FromTransaction: &CreateResponseService{
			ID:              fromTx.ID,
			WalletID:        fromTx.WalletID,
			CategoryID:      fromTx.CategoryID,
			Amount:          fromTx.Amount,
			Type:            fromTx.Type,
			TransactionDate: fromTx.TransactionDate,
			Note:            fromTx.Note,
			ImageURL:        s.resolveImageURL(ctx, fromTx.ImageURL),
			CreatedAt:       fromTx.CreatedAt,
			UpdatedAt:       fromTx.UpdatedAt,
		},
		ToTransaction: &CreateResponseService{
			ID:              toTx.ID,
			WalletID:        toTx.WalletID,
			CategoryID:      toTx.CategoryID,
			Amount:          toTx.Amount,
			Type:            toTx.Type,
			TransactionDate: toTx.TransactionDate,
			Note:            toTx.Note,
			ImageURL:        s.resolveImageURL(ctx, toTx.ImageURL),
			CreatedAt:       toTx.CreatedAt,
			UpdatedAt:       toTx.UpdatedAt,
		},
	}, nil
}
