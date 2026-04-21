package goals

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoParamController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoGoalController(ctx *gin.Context) {
	var param InfoParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoGoal(ctx, &InfoRequestService{ID: param.ID})
	if err != nil {
		if errors.Is(err, ErrGoalInvalidID) {
			_ = base.BadRequest(ctx, "goal-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrGoalNotFound) {
			_ = base.BadRequest(ctx, "goal-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "goal-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res)
}
