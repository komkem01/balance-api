package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	WalletID   *string `form:"wallet_id"`
	CategoryID *string `form:"category_id"`
	Type       *string `form:"type"`
	Page       int     `form:"page"`
	Size       int     `form:"size"`
}

func (c *Controller) ListTransactionController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, paginate, err := c.svc.ListTransaction(ctx, &ListRequestService{WalletID: req.WalletID, CategoryID: req.CategoryID, Type: req.Type, Page: req.Page, Size: req.Size})
	if err != nil {
		if errors.Is(err, ErrTransactionTypeInvalid) {
			_ = base.BadRequest(ctx, "transaction-type-invalid", gin.H{"field": "type", "reason": "invalid", "allowed": []string{"income", "expense"}})
			return
		}
		_ = base.InternalServerError(ctx, "transaction-list-failed", nil)
		return
	}
	_ = base.Paginate(ctx, res, paginate)
}
