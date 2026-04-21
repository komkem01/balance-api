package categories

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID  *string `json:"member_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Purpose   *string `json:"purpose"`
	IconName  string  `json:"icon_name"`
	ColorCode string  `json:"color_code"`
}

func (c *Controller) CreateCategoryController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateCategory(ctx, &CreateRequestService{MemberID: req.MemberID, Name: req.Name, Type: req.Type, Purpose: req.Purpose, IconName: req.IconName, ColorCode: req.ColorCode})
	if err != nil {
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
		if errors.Is(err, ErrCategoryInvalidMemberID) {
			_ = base.BadRequest(ctx, "category-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "category-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "category-created")
}
