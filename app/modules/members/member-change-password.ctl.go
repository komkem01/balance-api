package members

import (
	"errors"
	"strings"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ChangeMePasswordBodyController struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (c *Controller) ChangeMePasswordController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	var body ChangeMePasswordBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if strings.TrimSpace(body.ConfirmPassword) == "" || body.NewPassword != body.ConfirmPassword {
		_ = base.BadRequest(ctx, "member-password-confirmation-mismatch", gin.H{"field": "confirm_password", "reason": "mismatch"})
		return
	}

	res, err := c.svc.ChangeMePassword(ctx, &ChangeMePasswordRequestService{
		MemberID:        memberID,
		CurrentPassword: body.CurrentPassword,
		NewPassword:     body.NewPassword,
	})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberPasswordRequired) {
			_ = base.BadRequest(ctx, "member-password-required", nil)
			return
		}
		if errors.Is(err, ErrMemberInvalidCredentials) {
			_ = base.BadRequest(ctx, "member-current-password-invalid", gin.H{"field": "current_password", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberAccountNotFound) {
			_ = base.BadRequest(ctx, "member-account-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-password-change-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-password-changed")
}
