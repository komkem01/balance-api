package members

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GoogleOAuthStartController(ctx *gin.Context) {
	redirectURL, err := c.svc.GoogleOAuthStartURL(ctx)
	if err != nil {
		if errors.Is(err, ErrMemberGoogleOAuthDisabled) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": "google-oauth-disabled", "message": "Google OAuth is not enabled"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "google-oauth-start-failed", "message": "Google OAuth start failed"})
		return
	}

	ctx.Redirect(http.StatusFound, redirectURL)
}

func (c *Controller) GoogleOAuthCallbackController(ctx *gin.Context) {
	if oauthErr := ctx.Query("error"); oauthErr != "" {
		ctx.Redirect(http.StatusFound, c.svc.GoogleOAuthFailureRedirectURL("google-oauth-denied"))
		return
	}

	code := ctx.Query("code")
	state := ctx.Query("state")
	login, err := c.svc.GoogleOAuthCallback(ctx, code, state)
	if err != nil {
		reason := "google-oauth-failed"
		if errors.Is(err, ErrMemberGoogleInvalidState) {
			reason = "google-oauth-invalid-state"
		}
		if errors.Is(err, ErrMemberGoogleOAuthDisabled) {
			reason = "google-oauth-disabled"
		}
		ctx.Redirect(http.StatusFound, c.svc.GoogleOAuthFailureRedirectURL(reason))
		return
	}

	ctx.Redirect(http.StatusFound, c.svc.GoogleOAuthSuccessRedirectURL(login))
}
