package budgets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID   *string `form:"member_id"`
	CategoryID *string `form:"category_id"`
	Period     *string `form:"period"`
	Page       int     `form:"page"`
	Size       int     `form:"size"`
}

func (c *Controller) ListBudgetController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, paginate, err := c.svc.ListBudget(ctx, &ListRequestService{MemberID: req.MemberID, CategoryID: req.CategoryID, Period: req.Period, Page: req.Page, Size: req.Size})
	if err != nil {
		if errors.Is(err, ErrBudgetPeriodInvalid) {
			_ = base.BadRequest(ctx, "budget-period-invalid", gin.H{"field": "period", "reason": "invalid", "allowed": []string{"daily", "weekly", "monthly"}})
			return
		}
		_ = base.InternalServerError(ctx, "budget-list-failed", nil)
		return
	}

	_ = base.Paginate(ctx, res, paginate)
}
