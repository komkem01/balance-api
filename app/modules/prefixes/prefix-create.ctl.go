package prefixes

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	GenderID string `json:"gender_id" binding:"required"`
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

func (c *Controller) CreatePrefixController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.CreatePrefix(ctx, &CreateRequestService{
		GenderID: req.GenderID,
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, ErrPrefixNameRequired) {
			_ = base.BadRequest(ctx, "prefix-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}

		if errors.Is(err, ErrPrefixInvalidGenderID) {
			_ = base.BadRequest(ctx, "prefix-gender-id-invalid", gin.H{"field": "gender_id", "reason": "invalid"})
			return
		}

		if errors.Is(err, ErrPrefixAlreadyExists) {
			_ = base.BadRequest(ctx, "prefix-name-duplicate", gin.H{"field": "name", "reason": "duplicate"})
			return
		}

		_ = base.InternalServerError(ctx, "prefix-create-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "prefix-created")
}
