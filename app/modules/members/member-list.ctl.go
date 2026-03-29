package members

import (
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (c *Controller) ListMemberController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, paginate, err := c.svc.ListMember(ctx, &ListRequestService{Page: req.Page, Size: req.Size})
	if err != nil {
		_ = base.InternalServerError(ctx, "member-list-failed", nil)
		return
	}

	_ = base.Paginate(ctx, res, paginate)
}
