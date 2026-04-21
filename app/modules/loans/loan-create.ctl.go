package loans

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID         *string `json:"member_id"`
	Name             string  `json:"name"`
	Lender           string  `json:"lender"`
	TotalAmount      float64 `json:"total_amount"`
	RemainingBalance float64 `json:"remaining_balance"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	InterestRate     float64 `json:"interest_rate"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
}

func (c *Controller) CreateLoanController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.CreateLoan(ctx, &CreateRequestService{
		MemberID:         req.MemberID,
		Name:             req.Name,
		Lender:           req.Lender,
		TotalAmount:      req.TotalAmount,
		RemainingBalance: req.RemainingBalance,
		MonthlyPayment:   req.MonthlyPayment,
		InterestRate:     req.InterestRate,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
	})
	if err != nil {
		if errors.Is(err, ErrLoanNameRequired) {
			_ = base.BadRequest(ctx, "loan-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrLoanInvalidMemberID) {
			_ = base.BadRequest(ctx, "loan-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "loan-create-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "loan-created")
}
