package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type UpdateMeSettingsBodyController struct {
	PreferredCurrency *string `json:"preferred_currency"`
	PreferredLanguage *string `json:"preferred_language"`
	NotifyBudget      *bool   `json:"notify_budget"`
	NotifySecurity    *bool   `json:"notify_security"`
	NotifyWeekly      *bool   `json:"notify_weekly"`
}

type UpdateMeNotificationsBodyController struct {
	NotifyBudget   *bool `json:"notify_budget"`
	NotifySecurity *bool `json:"notify_security"`
	NotifyWeekly   *bool `json:"notify_weekly"`
}

func (c *Controller) InfoMeSettingsController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	res, err := c.svc.InfoMeSettings(ctx, memberID)
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-settings-info-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-settings-info")
}

func (c *Controller) UpdateMeSettingsController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	var body UpdateMeSettingsBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.UpdateMeSettings(ctx, &MeSettingsUpdateRequestService{
		MemberID:          memberID,
		PreferredCurrency: body.PreferredCurrency,
		PreferredLanguage: body.PreferredLanguage,
		NotifyBudget:      body.NotifyBudget,
		NotifySecurity:    body.NotifySecurity,
		NotifyWeekly:      body.NotifyWeekly,
	})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberInvalidCurrency) {
			_ = base.BadRequest(ctx, "member-settings-currency-invalid", gin.H{"field": "preferred_currency", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberInvalidLanguage) {
			_ = base.BadRequest(ctx, "member-settings-language-invalid", gin.H{"field": "preferred_language", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrMemberNoSettingsToUpdate) {
			_ = base.BadRequest(ctx, "member-settings-no-fields-to-update", nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-settings-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-settings-updated")
}

func (c *Controller) UpdateMeNotificationSettingsController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	var body UpdateMeNotificationsBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.UpdateMeSettings(ctx, &MeSettingsUpdateRequestService{
		MemberID:       memberID,
		NotifyBudget:   body.NotifyBudget,
		NotifySecurity: body.NotifySecurity,
		NotifyWeekly:   body.NotifyWeekly,
	})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberNoSettingsToUpdate) {
			_ = base.BadRequest(ctx, "member-settings-no-fields-to-update", nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-settings-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-settings-notifications-updated")
}
