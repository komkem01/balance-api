package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type RegisterRequestController struct {
	GenderID    *string `json:"gender_id"`
	PrefixID    *string `json:"prefix_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Phone       string  `json:"phone"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
}

func (c *Controller) RegisterMemberController(ctx *gin.Context) {
	var req RegisterRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.RegisterMember(ctx, &RegisterRequestService{
		GenderID:    req.GenderID,
		PrefixID:    req.PrefixID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Phone:       req.Phone,
		Username:    req.Username,
		Password:    req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrMemberInvalidGenderID) {
			_ = base.BadRequest(ctx, "member-gender-id-invalid", gin.H{"field": "gender_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberInvalidPrefixID) {
			_ = base.BadRequest(ctx, "member-prefix-id-invalid", gin.H{"field": "prefix_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberUsernameRequired) {
			_ = base.BadRequest(ctx, "member-username-required", gin.H{"field": "username", "reason": "required"})
			return
		}
		if errors.Is(err, ErrMemberPasswordRequired) {
			_ = base.BadRequest(ctx, "member-password-required", gin.H{"field": "password", "reason": "required"})
			return
		}
		if errors.Is(err, ErrMemberUsernameExists) {
			_ = base.BadRequest(ctx, "member-username-duplicate", gin.H{"field": "username", "reason": "duplicate"})
			return
		}
		_ = base.InternalServerError(ctx, "member-register-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-registered")
}
