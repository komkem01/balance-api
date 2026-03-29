package genders

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoGenderController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoGender(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrGenderInvalidID) {
			_ = base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		if errors.Is(err, ErrGenderNotFound) {
			_ = base.BadRequest(ctx, "gender-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "gender-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "gender-info")
}
