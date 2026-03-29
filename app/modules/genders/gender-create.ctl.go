package genders

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

type CreateResponseController struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Controller) CreateGenderController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateGender(ctx, &CreateRequestService{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, ErrGenderNameRequired) {
			_ = base.BadRequest(ctx, "gender-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}

		if errors.Is(err, ErrGenderAlreadyExists) {
			_ = base.BadRequest(ctx, "gender-name-duplicate", gin.H{"field": "name", "reason": "duplicate"})
			return
		}

		_ = base.InternalServerError(ctx, "gender-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "gender-created")
}
