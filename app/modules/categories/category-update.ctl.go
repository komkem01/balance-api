package categories

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	MemberID  *string `json:"member_id"`
	Name      *string `json:"name"`
	Type      *string `json:"type"`
	Purpose   *string `json:"purpose"`
	IconName  *string `json:"icon_name"`
	ColorCode *string `json:"color_code"`
}

func (c *Controller) UpdateCategoryController(ctx *gin.Context) {
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
	res, err := c.svc.UpdateCategory(ctx, &UpdateRequestService{ID: param.ID, MemberID: body.MemberID, Name: body.Name, Type: body.Type, Purpose: body.Purpose, IconName: body.IconName, ColorCode: body.ColorCode})
	if err != nil {
		if errors.Is(err, ErrCategoryInvalidID) {
			_ = base.BadRequest(ctx, "category-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrCategoryInvalidMemberID) {
			_ = base.BadRequest(ctx, "category-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrCategoryNameRequired) {
			_ = base.BadRequest(ctx, "category-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrCategoryTypeInvalid) {
			_ = base.BadRequest(ctx, "category-type-invalid", gin.H{"field": "type", "reason": "invalid", "allowed": []string{"income", "expense"}})
			return
		}
		if errors.Is(err, ErrCategoryPurposeInvalid) {
			_ = base.BadRequest(ctx, "category-purpose-invalid", gin.H{"field": "purpose", "reason": "invalid", "allowed": []string{"loan_repayment"}})
			return
		}
		if errors.Is(err, ErrCategoryNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		}
		if errors.Is(err, ErrCategoryNotFound) {
			_ = base.BadRequest(ctx, "category-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "category-update-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "category-updated")
}
