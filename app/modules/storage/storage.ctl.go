package storage

import (
	"errors"
	"strings"

	"balance/app/modules/net/httpx"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer trace.Tracer
	svc    *Service
	cli    *httpx.Client
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc, cli: httpx.NewClient()}
}

type UploadTransactionSlipResponse struct {
	ImageURL        string `json:"image_url"`
	DisplayImageURL string `json:"display_image_url"`
}

type GetTransactionSlipResponse struct {
	ImageURL        string `json:"image_url"`
	DisplayImageURL string `json:"display_image_url"`
}

func (c *Controller) UploadTransactionSlipController(ctx *gin.Context) {
	walletID := strings.TrimSpace(ctx.PostForm("wallet_id"))
	if walletID == "" {
		_ = base.BadRequest(ctx, "storage-wallet-id-required", gin.H{"field": "wallet_id", "reason": "required"})
		return
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		_ = base.BadRequest(ctx, "storage-image-required", gin.H{"field": "image", "reason": "required"})
		return
	}

	imageURL, err := c.svc.UploadTransactionSlip(ctx, walletID, fileHeader)
	if err != nil {
		if errors.Is(err, ErrStorageNotConfigured) {
			_ = base.InternalServerError(ctx, "storage-not-configured", nil)
			return
		}
		if errors.Is(err, ErrImageTooLarge) {
			_ = base.BadRequest(ctx, "storage-image-too-large", gin.H{"field": "image", "reason": "max-size-2mb"})
			return
		}
		if errors.Is(err, ErrImageInvalidType) || errors.Is(err, ErrImageInvalidContent) || errors.Is(err, ErrImageRequired) || errors.Is(err, ErrImageEmpty) {
			_ = base.BadRequest(ctx, "storage-image-invalid", gin.H{"field": "image", "reason": "invalid"})
			return
		}

		_ = base.InternalServerError(ctx, "storage-upload-failed", nil)
		return
	}

	_ = base.Success(ctx, &UploadTransactionSlipResponse{
		ImageURL:        imageURL,
		DisplayImageURL: c.svc.DisplayImageURL(ctx, imageURL),
	}, "storage-upload-success")
}

func (c *Controller) GetTransactionSlipController(ctx *gin.Context) {
	rawImageURL := ""
	if value, ok := ctx.Get("storage_image_url"); ok {
		if text, ok := value.(string); ok {
			rawImageURL = strings.TrimSpace(text)
		}
	}

	if rawImageURL == "" {
		rawImageURL = strings.TrimSpace(ctx.Query("image_url"))
	}

	if rawImageURL == "" {
		_ = base.BadRequest(ctx, "storage-image-url-required", gin.H{"field": "image_url", "reason": "required"})
		return
	}

	_ = base.Success(ctx, &GetTransactionSlipResponse{
		ImageURL:        rawImageURL,
		DisplayImageURL: c.svc.DisplayImageURL(ctx, rawImageURL),
	}, "storage-get-success")
}

