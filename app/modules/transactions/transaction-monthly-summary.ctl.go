package transactions

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type MonthlySummaryRequestController struct {
	MemberID   *string `form:"member_id"`
	WalletID   *string `form:"wallet_id"`
	CategoryID *string `form:"category_id"`
	StartDate  *string `form:"start_date"`
	EndDate    *string `form:"end_date"`
	Range      *string `form:"range"`
}

func (c *Controller) MonthlySummaryTransactionController(ctx *gin.Context) {
	var req MonthlySummaryRequestController
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.MonthlySummaryTransaction(ctx, &MonthlySummaryRequestService{
		MemberID:   req.MemberID,
		WalletID:   req.WalletID,
		CategoryID: req.CategoryID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Range:      req.Range,
	})
	if err != nil {
		if errors.Is(err, ErrTransactionDateInvalid) {
			_ = base.BadRequest(ctx, "transaction-date-invalid", gin.H{"reason": "invalid", "format": "2006-01-02"})
			return
		}
		if errors.Is(err, ErrTransactionRangeInvalid) {
			_ = base.BadRequest(ctx, "transaction-range-invalid", gin.H{"field": "range", "reason": "invalid", "allowed": []string{"1d", "1w", "1m", "1y", "all"}})
			return
		}

		_ = base.InternalServerError(ctx, "transaction-monthly-summary-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "transaction-monthly-summary")
}
