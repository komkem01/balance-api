package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type LoginRequestController struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Controller) LoginMemberController(ctx *gin.Context) {
	var req LoginRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.LoginMember(ctx, &LoginRequestService{Username: req.Username, Password: req.Password})
	if err != nil {
		if errors.Is(err, ErrMemberInvalidCredentials) {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-credentials"})
			return
		}
		_ = base.InternalServerError(ctx, "member-login-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-login-success")
}

type RefreshTokenRequestController struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *Controller) RefreshMemberTokenController(ctx *gin.Context) {
	var req RefreshTokenRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.RefreshMemberToken(ctx, &RefreshTokenRequestService{RefreshToken: req.RefreshToken})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-refresh-token"})
			return
		}
		_ = base.InternalServerError(ctx, "member-refresh-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-refresh-success")
}
