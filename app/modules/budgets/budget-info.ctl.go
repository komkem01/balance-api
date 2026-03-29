package budgets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoBudgetController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoBudget(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrBudgetInvalidID) {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrBudgetNotFound) {
			_ = base.BadRequest(ctx, "budget-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "budget-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "budget-info")
}
