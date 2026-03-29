package budgets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	MemberID   *string  `json:"member_id"`
	CategoryID *string  `json:"category_id"`
	Amount     *float64 `json:"amount"`
	Period     *string  `json:"period"`
	StartDate  *string  `json:"start_date"`
	EndDate    *string  `json:"end_date"`
}

func (c *Controller) UpdateBudgetController(ctx *gin.Context) {
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

	res, err := c.svc.UpdateBudget(ctx, &UpdateRequestService{ID: param.ID, MemberID: body.MemberID, CategoryID: body.CategoryID, Amount: body.Amount, Period: body.Period, StartDate: body.StartDate, EndDate: body.EndDate})
	if err != nil {
		if errors.Is(err, ErrBudgetInvalidID) {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrBudgetInvalidMemberID) {
			_ = base.BadRequest(ctx, "budget-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrBudgetInvalidCategoryID) {
			_ = base.BadRequest(ctx, "budget-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrBudgetPeriodInvalid) {
			_ = base.BadRequest(ctx, "budget-period-invalid", gin.H{"field": "period", "reason": "invalid", "allowed": []string{"daily", "weekly", "monthly"}})
			return
		}
		if errors.Is(err, ErrBudgetDateInvalid) {
			_ = base.BadRequest(ctx, "budget-date-invalid", gin.H{"reason": "invalid-date", "format": "2006-01-02"})
			return
		}
		if errors.Is(err, ErrBudgetNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		}
		if errors.Is(err, ErrBudgetNotFound) {
			_ = base.BadRequest(ctx, "budget-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "budget-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "budget-updated")
}
