package goals

import (
	"errors"

	"balance/app/modules/entities/ent"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	Name               *string  `json:"name"`
	TargetAmount       *float64 `json:"target_amount"`
	StartAmount        *float64 `json:"start_amount"`
	CurrentAmount      *float64 `json:"current_amount"`
	StartDate          *string  `json:"start_date"`
	TargetDate         *string  `json:"target_date"`
	Status             *string  `json:"status"`
	AutoTracking       *bool    `json:"auto_tracking"`
	TrackingSourceType *string  `json:"tracking_source_type"`
	TrackingSourceID   *string  `json:"tracking_source_id"`
	DepositWalletID    *string  `json:"deposit_wallet_id"`
}

func (c *Controller) UpdateGoalController(ctx *gin.Context) {
	var param UpdateParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var body UpdateBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var status *ent.GoalStatus
	if body.Status != nil {
		v := ent.GoalStatus(*body.Status)
		status = &v
	}

	var sourceType *ent.GoalTrackingSourceType
	if body.TrackingSourceType != nil {
		v := ent.GoalTrackingSourceType(*body.TrackingSourceType)
		sourceType = &v
	}

	res, err := c.svc.UpdateGoal(ctx, &UpdateRequestService{
		ID:                 param.ID,
		Name:               body.Name,
		TargetAmount:       body.TargetAmount,
		StartAmount:        body.StartAmount,
		CurrentAmount:      body.CurrentAmount,
		StartDate:          body.StartDate,
		TargetDate:         body.TargetDate,
		Status:             status,
		AutoTracking:       body.AutoTracking,
		TrackingSourceType: sourceType,
		TrackingSourceID:   body.TrackingSourceID,
		DepositWalletID:    body.DepositWalletID,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrGoalInvalidID):
			_ = base.BadRequest(ctx, "goal-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalNotFound):
			_ = base.BadRequest(ctx, "goal-not-found", nil)
			return
		case errors.Is(err, ErrGoalNameRequired):
			_ = base.BadRequest(ctx, "goal-name-required", gin.H{"field": "name", "reason": "required"})
			return
		case errors.Is(err, ErrGoalNoFieldsToUpdate):
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		case errors.Is(err, ErrGoalStatusInvalid):
			_ = base.BadRequest(ctx, "goal-status-invalid", gin.H{"field": "status", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalTargetAmountInvalid):
			_ = base.BadRequest(ctx, "goal-target-amount-invalid", gin.H{"field": "target_amount", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalStartDateInvalid):
			_ = base.BadRequest(ctx, "goal-start-date-invalid", gin.H{"field": "start_date", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalTargetDateInvalid):
			_ = base.BadRequest(ctx, "goal-target-date-invalid", gin.H{"field": "target_date", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalSourceTypeInvalid):
			_ = base.BadRequest(ctx, "goal-source-type-invalid", gin.H{"field": "tracking_source_type", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalSourceIDRequired):
			_ = base.BadRequest(ctx, "goal-source-id-required", gin.H{"field": "tracking_source_id", "reason": "required"})
			return
		case errors.Is(err, ErrGoalSourceIDInvalid):
			_ = base.BadRequest(ctx, "goal-source-id-invalid", gin.H{"field": "tracking_source_id", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalSourceMemberForbidden):
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-goal-source"})
			return
		case errors.Is(err, ErrGoalDepositWalletInvalid):
			_ = base.BadRequest(ctx, "goal-deposit-wallet-id-invalid", gin.H{"field": "deposit_wallet_id", "reason": "invalid"})
			return
		case errors.Is(err, ErrGoalDepositWalletForbidden):
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-goal-deposit-wallet"})
			return
		default:
			_ = base.InternalServerError(ctx, "goal-update-failed", nil)
			return
		}
	}

	_ = base.Success(ctx, res, "goal-updated")
}
