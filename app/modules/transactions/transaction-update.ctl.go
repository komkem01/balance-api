package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	WalletID        *string  `json:"wallet_id"`
	CategoryID      *string  `json:"category_id"`
	Amount          *float64 `json:"amount"`
	Type            *string  `json:"type"`
	TransactionDate *string  `json:"transaction_date"`
	Note            *string  `json:"note"`
	ImageURL        *string  `json:"image_url"`
}

func (c *Controller) UpdateTransactionController(ctx *gin.Context) {
	var param UpdateParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var body UpdateBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.UpdateTransaction(ctx, &UpdateRequestService{ID: param.ID, WalletID: body.WalletID, CategoryID: body.CategoryID, Amount: body.Amount, Type: body.Type, TransactionDate: body.TransactionDate, Note: body.Note, ImageURL: body.ImageURL})
	if err != nil {
		if errors.Is(err, ErrTransactionInvalidID) {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
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
		if errors.Is(err, ErrTransactionDateInvalid) {
			_ = base.BadRequest(ctx, "transaction-date-invalid", gin.H{"field": "transaction_date", "reason": "invalid", "format": "2006-01-02"})
			return
		}
		if errors.Is(err, ErrTransactionNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		}
		if errors.Is(err, ErrTransactionNotFound) {
			_ = base.BadRequest(ctx, "transaction-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "transaction-update-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "transaction-updated")
}
