package categories

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteCategoryController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	if err := c.svc.DeleteCategory(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrCategoryInvalidID) {
			_ = base.BadRequest(ctx, "category-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "category-delete-failed", nil)
		return
	}
	_ = base.Success(ctx, nil, "category-deleted")
}
