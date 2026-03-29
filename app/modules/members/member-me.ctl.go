package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

func resolveCurrentMemberID(ctx *gin.Context) string {
	if value, ok := ctx.Get("member_id"); ok {
		if text, ok := value.(string); ok && text != "" {
			return text
		}
	}
	return ""
}

func (c *Controller) InfoMeMemberController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	res, err := c.svc.InfoMeMember(ctx, &MeRequestService{MemberID: memberID})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-member-id"})
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-me-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-me")
}
