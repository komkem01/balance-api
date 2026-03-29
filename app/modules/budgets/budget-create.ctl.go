package budgets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID   *string `json:"member_id"`
	CategoryID *string `json:"category_id"`
	Amount     float64 `json:"amount"`
	Period     string  `json:"period"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
}

func (c *Controller) CreateBudgetController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateBudget(ctx, &CreateRequestService{MemberID: req.MemberID, CategoryID: req.CategoryID, Amount: req.Amount, Period: req.Period, StartDate: req.StartDate, EndDate: req.EndDate})
	if err != nil {
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
		_ = base.InternalServerError(ctx, "budget-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "budget-created")
}
