package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteTransactionController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteTransaction(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrTransactionInvalidID) {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionNotFound) {
			_ = base.BadRequest(ctx, "transaction-not-found", nil)
			return
		}
		if errors.Is(err, ErrTransactionInsufficientFunds) {
			_ = base.BadRequest(ctx, "transaction-insufficient-funds", gin.H{"field": "amount", "reason": "insufficient-wallet-balance"})
			return
		}
		_ = base.InternalServerError(ctx, "transaction-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "transaction-deleted")
}
