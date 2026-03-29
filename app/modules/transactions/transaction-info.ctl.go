package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoTransactionController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, err := c.svc.InfoTransaction(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrTransactionInvalidID) {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionNotFound) {
			_ = base.BadRequest(ctx, "transaction-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "transaction-info-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "transaction-info")
}
