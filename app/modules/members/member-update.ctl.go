package members

import (
	"errors"
	"time"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	GenderID    *string    `json:"gender_id"`
	PrefixID    *string    `json:"prefix_id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	DisplayName *string    `json:"display_name"`
	Phone       *string    `json:"phone"`
	LastLogin   *time.Time `json:"last_login"`
}

func (c *Controller) UpdateMemberController(ctx *gin.Context) {
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

	res, err := c.svc.UpdateMember(ctx, &UpdateRequestService{
		ID:          param.ID,
		GenderID:    body.GenderID,
		PrefixID:    body.PrefixID,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		DisplayName: body.DisplayName,
		Phone:       body.Phone,
		LastLogin:   body.LastLogin,
	})
	if err != nil {
		if errors.Is(err, ErrMemberInvalidID) {
			_ = base.BadRequest(ctx, "member-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberInvalidGenderID) {
			_ = base.BadRequest(ctx, "member-gender-id-invalid", gin.H{"field": "gender_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberInvalidPrefixID) {
			_ = base.BadRequest(ctx, "member-prefix-id-invalid", gin.H{"field": "prefix_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "member-no-fields-to-update", nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-updated")
}
