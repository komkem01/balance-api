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

func ownerTransactionCreateMiddleware(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberID, _ := ctx.Get("member_id")
		memberIDStr, _ := memberID.(string)

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

		if payload.WalletID == nil || strings.TrimSpace(*payload.WalletID) == "" {
			_ = base.BadRequest(ctx, "transaction-wallet-id-invalid", gin.H{"field": "wallet_id", "reason": "required"})
			ctx.Abort()
			return
		}

		wallet, err := mod.ENT.Svc.GetWalletByID(ctx, strings.TrimSpace(*payload.WalletID))
		if err != nil || wallet.MemberID == nil || wallet.MemberID.String() != memberIDStr {
			_ = base.Unauthorized(ctx, "unauthorized", gin.H{"reason": "forbidden-wallet"})
			ctx.Abort()
			return
		}

		if payload.CategoryID != nil && strings.TrimSpace(*payload.CategoryID) != "" {
			category, err := mod.ENT.Svc.GetCategoryByID(ctx, strings.TrimSpace(*payload.CategoryID))
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

		writeBody(ctx, body)
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
