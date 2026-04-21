package goals

import (
	"balance/app/modules/entities/ent"
	"balance/app/utils/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type ListRequestController struct {
	MemberID *string `form:"member_id"`
	Status   *string `form:"status"`
	Type     *string `form:"type"`
	Page     int     `form:"page"`
	Size     int     `form:"size"`
}

func (c *Controller) ListGoalController(ctx *gin.Context) {
	var req ListRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var status *ent.GoalStatus
	if req.Status != nil {
		v := ent.GoalStatus(strings.TrimSpace(*req.Status))
		if v != "" && !isValidGoalStatus(v) {
			_ = base.BadRequest(ctx, "goal-status-invalid", gin.H{"field": "status", "reason": "invalid"})
			return
		}
		if v != "" {
			status = &v
		}
	}

	var goalType *ent.GoalType
	if req.Type != nil {
		v := ent.GoalType(strings.TrimSpace(*req.Type))
		if v != "" && !isValidGoalType(v) {
			_ = base.BadRequest(ctx, "goal-type-invalid", gin.H{"field": "type", "reason": "invalid"})
			return
		}
		if v != "" {
			goalType = &v
		}
	}

	res, paginate, err := c.svc.ListGoal(ctx, &ListRequestService{MemberID: req.MemberID, Status: status, Type: goalType, Page: req.Page, Size: req.Size})
	if err != nil {
		_ = base.InternalServerError(ctx, "goal-list-failed", nil)
		return
	}

	_ = base.Paginate(ctx, res, paginate)
}
