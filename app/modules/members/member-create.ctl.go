package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	GenderID    *string `json:"gender_id"`
	PrefixID    *string `json:"prefix_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Phone       string  `json:"phone"`
}

func (c *Controller) CreateMemberController(ctx *gin.Context) {
	var req CreateRequestController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
		return
	}

	res, err := c.svc.CreateMember(ctx, &CreateRequestService{
		GenderID:    req.GenderID,
		PrefixID:    req.PrefixID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Phone:       req.Phone,
	})
	if err != nil {
		if errors.Is(err, ErrMemberInvalidGenderID) {
			_ = base.BadRequest(ctx, "member-gender-id-invalid", gin.H{"field": "gender_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberInvalidPrefixID) {
			_ = base.BadRequest(ctx, "member-prefix-id-invalid", gin.H{"field": "prefix_id", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-create-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-created")
}
