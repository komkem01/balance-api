package members

import (
	"errors"

	"balance/app/modules/storage"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetMeProfileImageController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	res, err := c.svc.InfoMeMember(ctx, &MeRequestService{MemberID: memberID})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-member-id"})
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		_ = base.InternalServerError(ctx, "member-profile-image-get-failed", nil)
		return
	}

	_ = base.Success(ctx, gin.H{"image_url": res.ProfileImageURL}, "member-profile-image-fetched")
}

func (c *Controller) UploadMeProfileImageController(ctx *gin.Context) {
	memberID := resolveCurrentMemberID(ctx)
	if memberID == "" {
		_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
		return
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		_ = base.BadRequest(ctx, "member-profile-image-required", gin.H{"field": "image", "reason": "required"})
		return
	}

	res, err := c.svc.UploadMeProfileImage(ctx, &UploadProfileImageRequestService{MemberID: memberID, Image: fileHeader})
	if err != nil {
		if errors.Is(err, ErrMemberUnauthorized) {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-member-id"})
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			_ = base.BadRequest(ctx, "member-not-found", nil)
			return
		}
		if errors.Is(err, ErrMemberStorageNotConfigured) {
			_ = base.InternalServerError(ctx, "member-profile-image-storage-not-configured", nil)
			return
		}
		if errors.Is(err, storage.ErrProfileImageTooLarge) {
			_ = base.BadRequest(ctx, "member-profile-image-too-large", gin.H{"field": "image", "reason": "max-size-10mb"})
			return
		}
		if errors.Is(err, storage.ErrImageInvalidType) || errors.Is(err, storage.ErrImageInvalidContent) || errors.Is(err, storage.ErrImageRequired) || errors.Is(err, storage.ErrImageEmpty) {
			_ = base.BadRequest(ctx, "member-profile-image-invalid", gin.H{"field": "image", "reason": "invalid"})
			return
		}
		_ = base.InternalServerError(ctx, "member-profile-image-upload-failed", nil)
		return
	}

	_ = base.Success(ctx, res, "member-profile-image-updated")
}
