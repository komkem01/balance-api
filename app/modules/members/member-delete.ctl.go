package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteMemberController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteMember(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrMemberInvalidID) {
			_ = base.BadRequest(ctx, "member-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "member-deleted")
}
