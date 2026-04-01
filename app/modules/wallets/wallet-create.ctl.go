package wallets

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID  *string `json:"member_id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	ColorCode string  `json:"color_code"`
	IsActive  *bool   `json:"is_active"`
}

func (c *Controller) CreateWalletController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateWallet(ctx, &CreateRequestService{MemberID: req.MemberID, Name: req.Name, Balance: req.Balance, Currency: req.Currency, ColorCode: req.ColorCode, IsActive: req.IsActive})
	if err != nil {
		if errors.Is(err, ErrWalletNameRequired) {
			_ = base.BadRequest(ctx, "wallet-name-required", gin.H{"field": "name", "reason": "required"})
			return
		}
		if errors.Is(err, ErrWalletInvalidMemberID) {
			_ = base.BadRequest(ctx, "wallet-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrWalletBalanceInvalid) {
			_ = base.BadRequest(ctx, "wallet-balance-invalid", gin.H{"field": "balance", "reason": "must-be-non-negative"})
			return
		}
		_ = base.InternalServerError(ctx, "wallet-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "wallet-created")
}
