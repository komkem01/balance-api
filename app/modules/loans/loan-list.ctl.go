package loans

import (
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID *string `form:"member_id"`
	Page     int     `form:"page"`
	Size     int     `form:"size"`
}

func (c *Controller) ListLoanController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, paginate, err := c.svc.ListLoan(ctx, &ListRequestService{MemberID: req.MemberID, Page: req.Page, Size: req.Size})
	if err != nil {
		_ = base.InternalServerError(ctx, "loan-list-failed", nil)
		return
	}

	_ = base.Paginate(ctx, res, paginate)
}
