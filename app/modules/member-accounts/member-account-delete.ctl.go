package memberaccounts

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type DeleteRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) DeleteMemberAccountController(ctx *gin.Context) {
	var req DeleteRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if err := c.svc.DeleteMemberAccount(ctx, &DeleteRequestService{ID: req.ID}); err != nil {
		if errors.Is(err, ErrMemberAccountInvalidID) {
			_ = base.BadRequest(ctx, "member-account-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-account-delete-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "member-account-deleted")
}
