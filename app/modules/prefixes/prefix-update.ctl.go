package prefixes

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

func (c *Controller) UpdatePrefixController(ctx *gin.Context) {
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

	res, err := c.svc.UpdatePrefix(ctx, &UpdateRequestService{
		ID:       param.ID,
		Name:     body.Name,
		IsActive: body.IsActive,
	})
	if err != nil {
		if errors.Is(err, ErrPrefixInvalidID) {
			_ = base.BadRequest(ctx, "prefix-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrPrefixNameRequired) {
			_ = base.BadRequest(ctx, "prefix-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrPrefixNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "prefix-no-fields-to-update", gin.H{"reason": "empty-update"})
			return
		}
		if errors.Is(err, ErrPrefixAlreadyExists) {
			_ = base.BadRequest(ctx, "prefix-name-duplicate", gin.H{"field": "name", "reason": "duplicate"})
			return
		}
		if errors.Is(err, ErrPrefixNotFound) {
			_ = base.BadRequest(ctx, "prefix-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "prefix-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "prefix-updated")
}
