package loans

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteParamController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteLoanController(ctx *gin.Context) {
	var param DeleteParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteLoan(ctx, &DeleteRequestService{ID: param.ID}); err != nil {
		if errors.Is(err, ErrLoanInvalidID) {
			_ = base.BadRequest(ctx, "loan-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrLoanNotFound) {
			_ = base.BadRequest(ctx, "loan-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "loan-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "loan-deleted")
}
