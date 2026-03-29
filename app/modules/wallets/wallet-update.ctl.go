package wallets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateParamController struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateBodyController struct {
	MemberID  *string  `json:"member_id"`
	Name      *string  `json:"name"`
	Balance   *float64 `json:"balance"`
	Currency  *string  `json:"currency"`
	ColorCode *string  `json:"color_code"`
	IsActive  *bool    `json:"is_active"`
}

func (c *Controller) UpdateWalletController(ctx *gin.Context) {
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
	res, err := c.svc.UpdateWallet(ctx, &UpdateRequestService{ID: param.ID, MemberID: body.MemberID, Name: body.Name, Balance: body.Balance, Currency: body.Currency, ColorCode: body.ColorCode, IsActive: body.IsActive})
	if err != nil {
		if errors.Is(err, ErrWalletInvalidID) {
			_ = base.BadRequest(ctx, "wallet-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrWalletInvalidMemberID) {
			_ = base.BadRequest(ctx, "wallet-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrWalletNameRequired) || errors.Is(err, ErrWalletNoFieldsToUpdate) {
			_ = base.BadRequest(ctx, "invalid-request", nil)
			return
		}
		if errors.Is(err, ErrWalletNotFound) {
			_ = base.BadRequest(ctx, "wallet-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "wallet-update-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "wallet-updated")
}
