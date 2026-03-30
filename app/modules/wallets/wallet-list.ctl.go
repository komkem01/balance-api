package wallets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID *string `form:"member_id"`
	IsActive *bool   `form:"is_active"`
	Page     int     `form:"page"`
	Size     int     `form:"size"`
}

func (c *Controller) ListWalletController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, paginate, err := c.svc.ListWallet(ctx, &ListRequestService{MemberID: req.MemberID, IsActive: req.IsActive, Page: req.Page, Size: req.Size})
	if err != nil {
		if errors.Is(err, ErrWalletInvalidMemberID) {
			_ = base.BadRequest(ctx, "wallet-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "wallet-list-failed", nil)
		return
	}
	_ = base.Paginate(ctx, res, paginate)
}
