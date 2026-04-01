package budgets

import (
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) RecalculateAllBudgetController(ctx *gin.Context) {
	res, err := c.svc.RecalculateAllBudgets(ctx)
	if err != nil {
		_ = base.InternalServerError(ctx, "budget-recalculate-all-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "budget-recalculate-all-done")
}
