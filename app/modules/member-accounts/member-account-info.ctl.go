package memberaccounts

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type InfoRequestController struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) InfoMemberAccountController(ctx *gin.Context) {
	var req InfoRequestController
	if err := ctx.ShouldBindUri(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.InfoMemberAccount(ctx, &InfoRequestService{ID: req.ID})
	if err != nil {
		if errors.Is(err, ErrMemberAccountInvalidID) {
			_ = base.BadRequest(ctx, "member-account-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberAccountNotFound) {
			_ = base.BadRequest(ctx, "member-account-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-account-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-account-info")
}
