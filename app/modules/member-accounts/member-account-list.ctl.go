package memberaccounts

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID *string `form:"member_id"`
	Page     int     `form:"page"`
	Size     int     `form:"size"`
}

func (c *Controller) ListMemberAccountController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, paginate, err := c.svc.ListMemberAccount(ctx, &ListRequestService{MemberID: req.MemberID, Page: req.Page, Size: req.Size})
	if err != nil {
		if errors.Is(err, ErrMemberAccountInvalidMemberID) {
			_ = base.BadRequest(ctx, "member-account-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-account-list-failed", nil)
		return
	}

	_ = base.Paginate(ctx, res, paginate)
}
