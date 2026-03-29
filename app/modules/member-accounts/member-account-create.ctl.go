package memberaccounts

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	MemberID *string `json:"member_id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
}

func (c *Controller) CreateMemberAccountController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateMemberAccount(ctx, &CreateRequestService{MemberID: req.MemberID, Username: req.Username, Password: req.Password})
	if err != nil {
		if errors.Is(err, ErrMemberAccountInvalidMemberID) {
			_ = base.BadRequest(ctx, "member-account-member-id-invalid", gin.H{"field": "member_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-account-create-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-account-created")
}
