package wallets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoWalletController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, err := c.svc.InfoWallet(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrWalletInvalidID) {
			_ = base.BadRequest(ctx, "wallet-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrWalletNotFound) {
			_ = base.BadRequest(ctx, "wallet-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "wallet-info-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "wallet-info")
}
