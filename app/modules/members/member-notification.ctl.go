package members

import (
	"errors"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListMeNotificationsQueryController struct {
	IncludeRead bool `form:"include_read"`
	Limit       int  `form:"limit"`
}

type SetMeNotificationReadBodyController struct {
	IsRead *bool `json:"is_read"`
}

func (c *Controller) ListMeNotificationsController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	var q ListMeNotificationsQueryController
	if err := ctx.ShouldBindQuery(&q); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.ListMeNotifications(ctx, &ListMeNotificationsRequestService{
		MemberID:    memberID,
		IncludeRead: q.IncludeRead,
		Limit:       q.Limit,
	})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-notification-list-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-notification-list")
}

func (c *Controller) SetMeNotificationReadController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	notificationID := ctx.Param("id")

	var body SetMeNotificationReadBodyController
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	isRead := true
	if body.IsRead != nil {
		isRead = *body.IsRead
	}

	res, err := c.svc.SetMeNotificationRead(ctx, &SetMeNotificationReadRequestService{
		MemberID:       memberID,
		NotificationID: notificationID,
		IsRead:         isRead,
	})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		if errors.Is(err, ErrMemberNotificationInvalidID) {
			_ = base.BadRequest(ctx, "member-notification-invalid-id", nil)
			return
		}
		if errors.Is(err, ErrMemberNotificationNotFound) {
			_ = base.BadRequest(ctx, "member-notification-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-notification-update-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-notification-updated")
}

func (c *Controller) MarkAllMeNotificationsReadController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	if err := c.svc.MarkAllMeNotificationsRead(ctx, memberID); err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-notification-update-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "member-notification-marked-all-read")
}

func (c *Controller) ClearMeNotificationsController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	if err := c.svc.ClearMeNotifications(ctx, memberID); err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-notification-clear-failed", nil)
		return
	}

	_ = base.Success(ctx, nil, "member-notification-cleared")
}
