package goals

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListEntriesParamController struct {
	ID string `uri:"id" binding:"required"`
}

type ListEntriesQueryController struct {
	MemberID *string `form:"member_id"`
}

func (c *Controller) ListGoalEntriesController(ctx *gin.Context) {
	var param ListEntriesParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var query ListEntriesQueryController
	if err := ctx.ShouldBindQuery(&query); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.ListGoalEntries(ctx, &ListEntriesRequestService{GoalID: param.ID, MemberID: query.MemberID})
	if err != nil {
		if errors.Is(err, ErrGoalInvalidID) {
			_ = base.BadRequest(ctx, "goal-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrGoalNotFound) {
			_ = base.BadRequest(ctx, "goal-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "goal-entry-list-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "goal-entry-listed")
}
