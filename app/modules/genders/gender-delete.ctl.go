package genders

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteGenderController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteGender(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrGenderInvalidID) {
			_ = base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		_ = base.InternalServerError(ctx, "gender-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "gender-deleted")
}
