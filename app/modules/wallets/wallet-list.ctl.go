package wallets

import (
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	IsActive *bool `form:"is_active"`
	Page     int   `form:"page"`
	Size     int   `form:"size"`
}

func (c *Controller) ListWalletController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, paginate, err := c.svc.ListWallet(ctx, &ListRequestService{IsActive: req.IsActive, Page: req.Page, Size: req.Size})
	if err != nil {
		_ = base.InternalServerError(ctx, "wallet-list-failed", nil)
		return
	}
	_ = base.Paginate(ctx, res, paginate)
}
