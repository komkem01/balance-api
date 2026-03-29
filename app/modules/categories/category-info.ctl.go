package categories

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoCategoryController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	res, err := c.svc.InfoCategory(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrCategoryInvalidID) {
			_ = base.BadRequest(ctx, "category-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrCategoryNotFound) {
			_ = base.BadRequest(ctx, "category-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "category-info-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "category-info")
}
