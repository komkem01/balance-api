package memberaccounts

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	MemberID *string `json:"member_id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (c *Controller) UpdateMemberAccountController(ctx *gin.Context) {
	var param UpdateParamController
	if err := ctx.ShouldBindUri(&param); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var body UpdateBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.UpdateMemberAccount(ctx, &UpdateRequestService{ID: param.ID, MemberID: body.MemberID, Username: body.Username, Password: body.Password})
	if err != nil {
		if errors.Is(err, ErrMemberAccountInvalidID) {
			_ = base.BadRequest(ctx, "member-account-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberAccountInvalidMemberID) {
			_ = base.BadRequest(ctx, "member-account-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberAccountNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "member-account-no-fields-to-update", nil)
			return
		}
		if errors.Is(err, ErrMemberAccountNotFound) {
			_ = base.BadRequest(ctx, "member-account-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-account-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-account-updated")
}
