package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type TransferRequestController struct {
	FromWalletID    string  `json:"from_wallet_id"`
	ToWalletID      string  `json:"to_wallet_id"`
	CategoryID      *string `json:"category_id"`
	Amount          float64 `json:"amount"`
	TransactionDate *string `json:"transaction_date"`
	Note            string  `json:"note"`
}

func (c *Controller) TransferBetweenWalletsController(ctx *gin.Context) {
	var req TransferRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.TransferBetweenWallets(ctx, &TransferRequestService{
		FromWalletID:    req.FromWalletID,
		ToWalletID:      req.ToWalletID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		TransactionDate: req.TransactionDate,
		Note:            req.Note,
	})
	if err != nil {
		if errors.Is(err, ErrTransactionTransferInvalidFromWalletID) {
			_ = base.BadRequest(ctx, "transaction-transfer-from-wallet-id-invalid", gin.H{"field": "from_wallet_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionTransferInvalidToWalletID) {
			_ = base.BadRequest(ctx, "transaction-transfer-to-wallet-id-invalid", gin.H{"field": "to_wallet_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionTransferSameWallet) {
			_ = base.BadRequest(ctx, "transaction-transfer-same-wallet", gin.H{"field": "to_wallet_id", "reason": "must-be-different"})
			return
		}
		if errors.Is(err, ErrTransactionInvalidCategoryID) {
			_ = base.BadRequest(ctx, "transaction-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionAmountInvalid) {
			_ = base.BadRequest(ctx, "transaction-amount-invalid", gin.H{"field": "amount", "reason": "must-be-greater-than-zero"})
			return
		}
		if errors.Is(err, ErrTransactionInsufficientFunds) {
			_ = base.BadRequest(ctx, "transaction-insufficient-funds", gin.H{"field": "amount", "reason": "insufficient-wallet-balance"})
			return
		}
		if errors.Is(err, ErrTransactionDateInvalid) {
			_ = base.BadRequest(ctx, "transaction-date-invalid", gin.H{"field": "transaction_date", "reason": "invalid", "format": "2006-01-02"})
			return
		}

		_ = base.InternalServerError(ctx, "transaction-transfer-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "transaction-transfer-created")
}
