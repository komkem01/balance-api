package budgets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteBudgetController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteBudget(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrBudgetInvalidID) {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "budget-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "budget-deleted")
}
