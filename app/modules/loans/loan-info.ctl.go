package loans

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoParamController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoLoanController(ctx *gin.Context) {
	var param InfoParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoLoan(ctx, &InfoRequestService{ID: param.ID})
	if err != nil {
		if errors.Is(err, ErrLoanInvalidID) {
			_ = base.BadRequest(ctx, "loan-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrLoanNotFound) {
			_ = base.BadRequest(ctx, "loan-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "loan-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res)
}
