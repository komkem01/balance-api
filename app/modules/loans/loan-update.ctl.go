package loans

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	Name             *string  `json:"name"`
	Lender           *string  `json:"lender"`
	TotalAmount      *float64 `json:"total_amount"`
	RemainingBalance *float64 `json:"remaining_balance"`
	MonthlyPayment   *float64 `json:"monthly_payment"`
	InterestRate     *float64 `json:"interest_rate"`
	StartDate        *string  `json:"start_date"`
	EndDate          *string  `json:"end_date"`
}

func (c *Controller) UpdateLoanController(ctx *gin.Context) {
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

	res, err := c.svc.UpdateLoan(ctx, &UpdateRequestService{
		ID:               param.ID,
		Name:             body.Name,
		Lender:           body.Lender,
		TotalAmount:      body.TotalAmount,
		RemainingBalance: body.RemainingBalance,
		MonthlyPayment:   body.MonthlyPayment,
		InterestRate:     body.InterestRate,
		StartDate:        body.StartDate,
		EndDate:          body.EndDate,
	})
	if err != nil {
		if errors.Is(err, ErrLoanInvalidID) {
			_ = base.BadRequest(ctx, "loan-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrLoanNotFound) {
			_ = base.BadRequest(ctx, "loan-not-found", nil)
			return
		}
		if errors.Is(err, ErrLoanNameRequired) {
			_ = base.BadRequest(ctx, "loan-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrLoanNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		}
		_ = base.InternalServerError(ctx, "loan-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "loan-updated")
}
