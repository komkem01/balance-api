package wallets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteWalletController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	if err := c.svc.DeleteWallet(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrWalletInvalidID) {
			_ = base.BadRequest(ctx, "wallet-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "wallet-delete-failed", nil)
		return
	}
	_ = base.Success(ctx, nil, "wallet-deleted")
}
