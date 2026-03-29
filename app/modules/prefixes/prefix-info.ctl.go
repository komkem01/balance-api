package prefixes

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoPrefixController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoPrefix(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrPrefixInvalidID) {
			_ = base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		if errors.Is(err, ErrPrefixNotFound) {
			_ = base.BadRequest(ctx, "prefix-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "prefix-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "prefix-info")
}
