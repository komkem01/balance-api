package goals

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteParamController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteGoalController(ctx *gin.Context) {
	var param DeleteParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteGoal(ctx, &DeleteRequestService{ID: param.ID}); err != nil {
		if errors.Is(err, ErrGoalInvalidID) {
			_ = base.BadRequest(ctx, "goal-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrGoalNotFound) {
			_ = base.BadRequest(ctx, "goal-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "goal-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "goal-deleted")
}
