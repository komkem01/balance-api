package goals

import (
	"errors"

	"balance/app/modules/entities/ent"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID           *string `json:"member_id"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	TargetAmount       float64 `json:"target_amount"`
	StartAmount        float64 `json:"start_amount"`
	CurrentAmount      float64 `json:"current_amount"`
	StartDate          *string `json:"start_date"`
	TargetDate         *string `json:"target_date"`
	Status             string  `json:"status"`
	AutoTracking       *bool   `json:"auto_tracking"`
	TrackingSourceType *string `json:"tracking_source_type"`
	TrackingSourceID   *string `json:"tracking_source_id"`
}

func (c *Controller) CreateGoalController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var sourceType *ent.GoalTrackingSourceType
	if req.TrackingSourceType != nil {
		v := ent.GoalTrackingSourceType(*req.TrackingSourceType)
		sourceType = &v
	}

	res, err := c.svc.CreateGoal(ctx, &CreateRequestService{
		MemberID:           req.MemberID,
		Name:               req.Name,
		Type:               ent.GoalType(req.Type),
		TargetAmount:       req.TargetAmount,
		StartAmount:        req.StartAmount,
		CurrentAmount:      req.CurrentAmount,
		StartDate:          req.StartDate,
		TargetDate:         req.TargetDate,
		Status:             ent.GoalStatus(req.Status),
		AutoTracking:       req.AutoTracking,
		TrackingSourceType: sourceType,
		TrackingSourceID:   req.TrackingSourceID,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrGoalNameRequired):
			_ = base.BadRequest(ctx, "goal-name-required", gin.H{"field": "name", "reason": "required"})
			return
		case errors.Is(err, ErrGoalTypeInvalid):
			_ = base.BadRequest(ctx, "goal-type-invalid", gin.H{"field": "type", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalStatusInvalid):
			_ = base.BadRequest(ctx, "goal-status-invalid", gin.H{"field": "status", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalTargetAmountInvalid):
			_ = base.BadRequest(ctx, "goal-target-amount-invalid", gin.H{"field": "target_amount", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalInvalidMemberID):
			_ = base.BadRequest(ctx, "goal-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalSourceTypeInvalid):
			_ = base.BadRequest(ctx, "goal-source-type-invalid", gin.H{"field": "tracking_source_type", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalSourceIDRequired):
			_ = base.BadRequest(ctx, "goal-source-id-required", gin.H{"field": "tracking_source_id", "reason": "required"})
			return
		case errors.Is(err, ErrGoalSourceMemberForbidden):
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-goal-source"})
			return
		default:
			_ = base.InternalServerError(ctx, "goal-create-failed", nil)
			return
		}
	}

	_ = base.Success(ctx, res, "goal-created")
}
