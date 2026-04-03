package transactions

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequestController struct {
	WalletID        *string `json:"wallet_id"`
	CategoryID      *string `json:"category_id"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	TransactionDate *string `json:"transaction_date"`
	Note            string  `json:"note"`
	ImageURL        string  `json:"image_url"`
	Image           *multipart.FileHeader
}

func ptrIfNotBlank(value string) *string {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}
	return &v
}

func (c *Controller) CreateTransactionController(ctx *gin.Context) {
	var req CreateRequestController
	if strings.HasPrefix(strings.ToLower(ctx.GetHeader("Content-Type")), "multipart/form-data") {
		amount, err := strconv.ParseFloat(strings.TrimSpace(ctx.PostForm("amount")), 64)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-amount"})
			return
		}

		typeValue := strings.TrimSpace(ctx.PostForm("type"))
		if typeValue == "" {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "type-required"})
			return
		}

		req = CreateRequestController{
			WalletID:        ptrIfNotBlank(ctx.PostForm("wallet_id")),
			CategoryID:      ptrIfNotBlank(ctx.PostForm("category_id")),
			Amount:          amount,
			Type:            typeValue,
			TransactionDate: ptrIfNotBlank(ctx.PostForm("transaction_date")),
			Note:            strings.TrimSpace(ctx.PostForm("note")),
			ImageURL:        strings.TrimSpace(ctx.PostForm("image_url")),
		}

		fileHeader, err := ctx.FormFile("image")
		if err == nil {
			req.Image = fileHeader
		}
	} else {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			return
		}
	}

	res, err := c.svc.CreateTransaction(ctx, &CreateRequestService{
		WalletID:        req.WalletID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		TransactionDate: req.TransactionDate,
		Note:            req.Note,
		ImageURL:        req.ImageURL,
		Image:           req.Image,
	})
	if err != nil {
		if errors.Is(err, ErrTransactionInvalidWalletID) {
			_ = base.BadRequest(ctx, "transaction-wallet-id-invalid", gin.H{"field": "wallet_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionInvalidCategoryID) {
			_ = base.BadRequest(ctx, "transaction-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionTypeInvalid) {
			_ = base.BadRequest(ctx, "transaction-type-invalid", gin.H{"field": "type", "reason": "invalid", "allowed": []string{"income", "expense"}})
			return
		}
		if errors.Is(err, ErrTransactionAmountInvalid) {
			_ = base.BadRequest(ctx, "transaction-amount-invalid", gin.H{"field": "amount", "reason": "must-be-non-negative"})
			return
		}
		if errors.Is(err, ErrTransactionInsufficientFunds) {
			_ = base.BadRequest(ctx, "transaction-insufficient-funds", gin.H{"field": "amount", "reason": "insufficient-wallet-balance"})
			return
		}
		if errors.Is(err, ErrTransactionDateInvalid) {
			_ = base.BadRequest(ctx, "transaction-date-invalid", gin.H{"field": "transaction_date", "reason": "invalid", "format": "2006-01-02"})
			return
		}
		if errors.Is(err, ErrTransactionImageInvalid) {
			_ = base.BadRequest(ctx, "transaction-image-invalid", gin.H{"field": "image", "reason": "invalid"})
			return
		}
		if errors.Is(err, ErrTransactionImageTooLarge) {
			_ = base.BadRequest(ctx, "transaction-image-too-large", gin.H{"field": "image", "reason": "max-size-10mb"})
			return
		}
		if errors.Is(err, ErrTransactionImageUploadFailed) {
			_ = base.InternalServerError(ctx, "transaction-image-upload-failed", nil)
			return
		}
		_ = base.InternalServerError(ctx, "transaction-create-failed", nil)
		return
	}
	_ = base.Success(ctx, res, "transaction-created")
}
