package genders

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

func (c *Controller) UpdateGenderController(ctx *gin.Context) {
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

	res, err := c.svc.UpdateGender(ctx, &UpdateRequestService{
		ID:       param.ID,
		Name:     body.Name,
		IsActive: body.IsActive,
	})
	if err != nil {
		if errors.Is(err, ErrGenderInvalidID) {
			_ = base.BadRequest(ctx, "gender-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrGenderNameRequired) {
			_ = base.BadRequest(ctx, "gender-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrGenderNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "gender-no-fields-to-update", gin.H{"reason": "empty-update"})
			return
		}
		if errors.Is(err, ErrGenderAlreadyExists) {
			_ = base.BadRequest(ctx, "gender-name-duplicate", gin.H{"field": "name", "reason": "duplicate"})
			return
		}
		if errors.Is(err, ErrGenderNotFound) {
			_ = base.BadRequest(ctx, "gender-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "gender-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "gender-updated")
}
