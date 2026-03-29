package categories

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID *string `form:"member_id"`
	Type     *string `form:"type"`
	Page     int     `form:"page"`
	Size     int     `form:"size"`
}

func (c *Controller) ListCategoryController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, paginate, err := c.svc.ListCategory(ctx, &ListRequestService{MemberID: req.MemberID, Type: req.Type, Page: req.Page, Size: req.Size})
	if err != nil {
		if errors.Is(err, ErrCategoryTypeInvalid) {
			_ = base.BadRequest(ctx, "category-type-invalid", gin.H{"field": "type", "reason": "invalid", "allowed": []string{"income", "expense"}})
			return
		}
		if errors.Is(err, ErrCategoryInvalidMemberID) {
			_ = base.BadRequest(ctx, "category-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "category-list-failed", nil)
		return
	}
	_ = base.Paginate(ctx, res, paginate)
}
