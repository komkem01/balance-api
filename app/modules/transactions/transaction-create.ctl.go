package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	WalletID        *string `json:"wallet_id"`
	CategoryID      *string `json:"category_id"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	TransactionDate *string `json:"transaction_date"`
	Note            string  `json:"note"`
	ImageURL        string  `json:"image_url"`
}

func (c *Controller) CreateTransactionController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateTransaction(ctx, &CreateRequestService{
		WalletID:        req.WalletID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		TransactionDate: req.TransactionDate,
		Note:            req.Note,
		ImageURL:        req.ImageURL,
	})
	if err != nil {
		if errors.Is(err, ErrTransactionInvalidWalletID) {
			_ = base.BadRequest(ctx, "transaction-wallet-id-invalid", gin.H{"field": "wallet_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionInvalidCategoryID) {
			_ = base.BadRequest(ctx, "transaction-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionTypeInvalid) {
			_ = base.BadRequest(ctx, "transaction-type-invalid", gin.H{"field": "type", "reason": "invalid", "allowed": []string{"income", "expense"}})
			return
		}
		if errors.Is(err, ErrTransactionAmountInvalid) {
			_ = base.BadRequest(ctx, "transaction-amount-invalid", gin.H{"field": "amount", "reason": "must-be-non-negative"})
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
		_ = base.InternalServerError(ctx, "transaction-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "transaction-created")
}
