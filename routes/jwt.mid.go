package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"balance/app/modules"
	"balance/app/utils/authx"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
)

func memberJWTMiddleware(mod *modules.Modules) gin.HandlerFunc {
	secret := strings.TrimSpace(mod.Conf.Svc.Config().AppKey)
	if secret == "" {
		secret = "secret"
	}

	return func(ctx *gin.Context) {
		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-bearer-token"})
			ctx.Abort()
			return
		}

		token := strings.TrimSpace(authorization[7:])
		claims, err := authx.ParseMemberToken(secret, token)
		if err != nil {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "invalid-token"})
			ctx.Abort()
			return
		}

		ctx.Set("member_id", claims.MemberID)
		ctx.Next()
	}
}

func requireMemberJWT(mod *modules.Modules) gin.HandlerFunc {
	return memberJWTMiddleware(mod)
}

func readAndRestoreBody(ctx *gin.Context) ([]byte, error) {
	if ctx.Request.Body == nil {
		return []byte("{}"), nil
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func writeBody(ctx *gin.Context, body []byte) {
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	ctx.Request.ContentLength = int64(len(body))
}

func resolveMemberIDFromContext(ctx *gin.Context) string {
	memberID, ok := ctx.Get("member_id")
	if !ok {
		return ""
	}
	memberIDStr, ok := memberID.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(memberIDStr)
}

func upsertMemberIDBodyMiddleware(force bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		body, err := readAndRestoreBody(ctx)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		trimmed := strings.TrimSpace(string(body))
		if trimmed == "" {
			if !force {
				writeBody(ctx, body)
				ctx.Next()
				return
			}
			body = []byte("{}")
		}

		payload := map[string]any{}
		if err := json.Unmarshal(body, &payload); err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		if force {
			payload["member_id"] = memberIDStr
		} else {
			if _, exists := payload["member_id"]; exists {
				payload["member_id"] = memberIDStr
			}
		}

		newBody, err := json.Marshal(payload)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		writeBody(ctx, newBody)
		ctx.Next()
	}
}

func forceMemberIDBodyMiddleware() gin.HandlerFunc {
	return upsertMemberIDBodyMiddleware(true)
}

func sanitizeMemberIDBodyMiddleware() gin.HandlerFunc {
	return upsertMemberIDBodyMiddleware(false)
}

func forceMemberIDQueryMiddleware(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		query := ctx.Request.URL.Query()
		query.Set(key, memberIDStr)
		ctx.Request.URL.RawQuery = query.Encode()

		ctx.Next()
	}
}

func ownerMemberAccountByParamMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		id := strings.TrimSpace(ctx.Param("id"))
		item, err := mod.ENT.Svc.GetMemberAccountByID(ctx, id)
		if err != nil {
			_ = base.BadRequest(ctx, "member-account-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if item.MemberID == nil || item.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-member-account"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerWalletByParamMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		id := strings.TrimSpace(ctx.Param("id"))
		item, err := mod.ENT.Svc.GetWalletByID(ctx, id)
		if err != nil {
			_ = base.BadRequest(ctx, "wallet-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if item.MemberID == nil || item.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-wallet"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerCategoryByParamMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		id := strings.TrimSpace(ctx.Param("id"))
		item, err := mod.ENT.Svc.GetCategoryByID(ctx, id)
		if err != nil {
			_ = base.BadRequest(ctx, "category-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if item.MemberID == nil || item.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-category"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerTransactionCreateMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		walletID := ""
		categoryID := ""
		contentType := strings.ToLower(strings.TrimSpace(ctx.GetHeader("Content-Type")))

		if strings.HasPrefix(contentType, "multipart/form-data") {
			walletID = strings.TrimSpace(ctx.PostForm("wallet_id"))
			categoryID = strings.TrimSpace(ctx.PostForm("category_id"))
		} else {
			body, err := readAndRestoreBody(ctx)
			if err != nil {
				_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
				ctx.Abort()
				return
			}

			var payload struct {
				WalletID   *string `json:"wallet_id"`
				CategoryID *string `json:"category_id"`
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
				ctx.Abort()
				return
			}

			if payload.WalletID != nil {
				walletID = strings.TrimSpace(*payload.WalletID)
			}
			if payload.CategoryID != nil {
				categoryID = strings.TrimSpace(*payload.CategoryID)
			}

			writeBody(ctx, body)
		}

		if walletID == "" {
			_ = base.BadRequest(ctx, "transaction-wallet-id-invalid", gin.H{"field": "wallet_id", "reason": "required"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, walletID)
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-wallet"})
			ctx.Abort()
			return
		}

		if categoryID != "" {
			category, err := mod.ENT.Svc.GetCategoryByID(ctx, categoryID)
			if err != nil {
				_ = base.BadRequest(ctx, "transaction-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
				ctx.Abort()
				return
			}
			if category.MemberID != nil && category.MemberID.String() != memberIDStr {
				_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-category"})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}

func ownerStorageUploadMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		walletID := strings.TrimSpace(ctx.PostForm("wallet_id"))
		if walletID == "" {
			_ = base.BadRequest(ctx, "storage-wallet-id-required", gin.H{"field": "wallet_id", "reason": "required"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, walletID)
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-wallet"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerStorageReadMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "missing-member-id"})
			ctx.Abort()
			return
		}

		transactionID := strings.TrimSpace(ctx.Param("id"))
		if transactionID == "" {
			_ = base.BadRequest(ctx, "storage-transaction-id-required", gin.H{"field": "id", "reason": "required"})
			ctx.Abort()
			return
		}

		transaction, err := mod.ENT.Svc.GetTransactionByID(ctx, transactionID)
		if err != nil || transaction.WalletID == nil {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, transaction.WalletID.String())
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-transaction"})
			ctx.Abort()
			return
		}

		imageURL := strings.TrimSpace(transaction.ImageURL)
		if imageURL == "" {
			_ = base.BadRequest(ctx, "storage-image-url-required", gin.H{"field": "image_url", "reason": "required"})
			ctx.Abort()
			return
		}

		ctx.Set("storage_image_url", imageURL)
		ctx.Next()
	}
}

func ownerBudgetCreateMiddleware(_ *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		body, err := readAndRestoreBody(ctx)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		payload := map[string]any{}
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
				ctx.Abort()
				return
			}
		}

		if raw, ok := payload["member_id"]; ok {
			if text, ok := raw.(string); ok && strings.TrimSpace(text) != "" && strings.TrimSpace(text) != memberIDStr {
				_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-member"})
				ctx.Abort()
				return
			}
		}

		payload["member_id"] = memberIDStr
		newBody, err := json.Marshal(payload)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		writeBody(ctx, newBody)
		ctx.Next()
	}
}

func ownerTransactionUpdateMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		transactionID := strings.TrimSpace(ctx.Param("id"))
		transaction, err := mod.ENT.Svc.GetTransactionByID(ctx, transactionID)
		if err != nil {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if transaction.WalletID == nil {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-transaction"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, transaction.WalletID.String())
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-transaction"})
			ctx.Abort()
			return
		}

		body, err := readAndRestoreBody(ctx)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		payload := map[string]any{}
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
				ctx.Abort()
				return
			}
		}

		if raw, ok := payload["wallet_id"]; ok {
			if text, ok := raw.(string); ok && strings.TrimSpace(text) != "" {
				newWallet, err := mod.ENT.Svc.GetWalletByID(ctx, strings.TrimSpace(text))
				if err != nil || newWallet.MemberID == nil || newWallet.MemberID.String() != memberIDStr {
					_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-wallet"})
					ctx.Abort()
					return
				}
			}
		}

		if raw, ok := payload["category_id"]; ok {
			if text, ok := raw.(string); ok && strings.TrimSpace(text) != "" {
				category, err := mod.ENT.Svc.GetCategoryByID(ctx, strings.TrimSpace(text))
				if err != nil {
					_ = base.BadRequest(ctx, "transaction-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
					ctx.Abort()
					return
				}
				if category.MemberID != nil && category.MemberID.String() != memberIDStr {
					_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-category"})
					ctx.Abort()
					return
				}
			}
		}

		writeBody(ctx, body)
		ctx.Next()
	}
}

func ownerTransactionDeleteMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		transactionID := strings.TrimSpace(ctx.Param("id"))
		transaction, err := mod.ENT.Svc.GetTransactionByID(ctx, transactionID)
		if err != nil || transaction.WalletID == nil {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, transaction.WalletID.String())
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-transaction"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerTransactionReadMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		transactionID := strings.TrimSpace(ctx.Param("id"))
		transaction, err := mod.ENT.Svc.GetTransactionByID(ctx, transactionID)
		if err != nil || transaction.WalletID == nil {
			_ = base.BadRequest(ctx, "transaction-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, transaction.WalletID.String())
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-transaction"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerBudgetUpdateMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		budgetID := strings.TrimSpace(ctx.Param("id"))
		budget, err := mod.ENT.Svc.GetBudgetByID(ctx, budgetID)
		if err != nil {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if budget.MemberID == nil || budget.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-budget"})
			ctx.Abort()
			return
		}

		body, err := readAndRestoreBody(ctx)
		if err != nil {
			_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
			ctx.Abort()
			return
		}

		payload := map[string]any{}
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				_ = base.BadRequest(ctx, "invalid-request", gin.H{"reason": "invalid-body"})
				ctx.Abort()
				return
			}
		}

		if raw, ok := payload["member_id"]; ok {
			if text, ok := raw.(string); ok {
				trimmed := strings.TrimSpace(text)
				if trimmed == "" || trimmed != memberIDStr {
					_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-member"})
					ctx.Abort()
					return
				}
			}
		}

		if raw, ok := payload["category_id"]; ok {
			if text, ok := raw.(string); ok && strings.TrimSpace(text) != "" {
				category, err := mod.ENT.Svc.GetCategoryByID(ctx, strings.TrimSpace(text))
				if err != nil {
					_ = base.BadRequest(ctx, "budget-category-id-invalid", gin.H{"field": "category_id", "reason": "invalid"})
					ctx.Abort()
					return
				}
				if category.MemberID != nil && category.MemberID.String() != memberIDStr {
					_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-category"})
					ctx.Abort()
					return
				}
			}
		}

		writeBody(ctx, body)
		ctx.Next()
	}
}

func ownerBudgetDeleteMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		budgetID := strings.TrimSpace(ctx.Param("id"))
		budget, err := mod.ENT.Svc.GetBudgetByID(ctx, budgetID)
		if err != nil {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if budget.MemberID == nil || budget.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-budget"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ownerBudgetReadMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

		budgetID := strings.TrimSpace(ctx.Param("id"))
		budget, err := mod.ENT.Svc.GetBudgetByID(ctx, budgetID)
		if err != nil {
			_ = base.BadRequest(ctx, "budget-invalid-id", gin.H{"field": "id", "reason": "invalid"})
			ctx.Abort()
			return
		}

		if budget.MemberID == nil || budget.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-budget"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
