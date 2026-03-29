package routes

import (
	"fmt"
	"net/http"

	"balance/app/modules"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

func api(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/example/:id", mod.Example.Ctl.Get)
	r.GET("/example-http", mod.Example.Ctl.GetHttpReq)
	r.POST("/example", mod.Example.Ctl.Create)
}

func apiSystem(r *gin.RouterGroup, mod *modules.Modules) {
	systems := r.Group("/systems")
	{
		genders := systems.Group("/genders")
		{
			genders.GET("", mod.Gender.Ctl.ListGenderController)
			genders.POST("", mod.Gender.Ctl.CreateGenderController)
			genders.GET("/:id", mod.Gender.Ctl.InfoGenderController)
			genders.PATCH("/:id", mod.Gender.Ctl.UpdateGenderController)
			genders.DELETE("/:id", mod.Gender.Ctl.DeleteGenderController)
		}
		prefixes := systems.Group("/prefixes")
		{
			prefixes.GET("", mod.Prefix.Ctl.ListPrefixController)
			prefixes.POST("", mod.Prefix.Ctl.CreatePrefixController)
			prefixes.GET("/:id", mod.Prefix.Ctl.InfoPrefixController)
			prefixes.PATCH("/:id", mod.Prefix.Ctl.UpdatePrefixController)
			prefixes.DELETE("/:id", mod.Prefix.Ctl.DeletePrefixController)
		}
	}
}

func apiMember(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/me", requireMemberJWT(mod), mod.Member.Ctl.InfoMeMemberController)

	members := r.Group("/members")
	{
		members.GET("", mod.Member.Ctl.ListMemberController)
		members.POST("", mod.Member.Ctl.CreateMemberController)
		members.GET("/:id", mod.Member.Ctl.InfoMemberController)
		members.PATCH("/:id", mod.Member.Ctl.UpdateMemberController)
		members.DELETE("/:id", mod.Member.Ctl.DeleteMemberController)
	}
	memberAccounts := r.Group("/member-accounts")
	{
		memberAccounts.GET("", mod.MemberAccount.Ctl.ListMemberAccountController)
		memberAccounts.POST("", mod.MemberAccount.Ctl.CreateMemberAccountController)
		memberAccounts.GET("/:id", mod.MemberAccount.Ctl.InfoMemberAccountController)
		memberAccounts.PATCH("/:id", mod.MemberAccount.Ctl.UpdateMemberAccountController)
		memberAccounts.DELETE("/:id", mod.MemberAccount.Ctl.DeleteMemberAccountController)
	}
}

func apiBalance(r *gin.RouterGroup, mod *modules.Modules) {
	Balances := r.Group("/balances")
	{
		wallets := Balances.Group("/wallets")
		{
			wallets.GET("", mod.Wallet.Ctl.ListWalletController)
			wallets.POST("", mod.Wallet.Ctl.CreateWalletController)
			wallets.GET("/:id", mod.Wallet.Ctl.InfoWalletController)
			wallets.PATCH("/:id", mod.Wallet.Ctl.UpdateWalletController)
			wallets.DELETE("/:id", mod.Wallet.Ctl.DeleteWalletController)
		}

		categories := Balances.Group("/categories")
		{
			categories.GET("", mod.Category.Ctl.ListCategoryController)
			categories.POST("", mod.Category.Ctl.CreateCategoryController)
			categories.GET("/:id", mod.Category.Ctl.InfoCategoryController)
			categories.PATCH("/:id", mod.Category.Ctl.UpdateCategoryController)
			categories.DELETE("/:id", mod.Category.Ctl.DeleteCategoryController)
		}

		transactions := Balances.Group("/transactions")
		{
			transactions.GET("", mod.Transaction.Ctl.ListTransactionController)
			transactions.POST("", requireMemberJWT(mod), ownerTransactionCreateMiddleware(mod), mod.Transaction.Ctl.CreateTransactionController)
			transactions.GET("/:id", requireMemberJWT(mod), ownerTransactionReadMiddleware(mod), mod.Transaction.Ctl.InfoTransactionController)
			transactions.PATCH("/:id", requireMemberJWT(mod), ownerTransactionUpdateMiddleware(mod), mod.Transaction.Ctl.UpdateTransactionController)
			transactions.DELETE("/:id", requireMemberJWT(mod), ownerTransactionDeleteMiddleware(mod), mod.Transaction.Ctl.DeleteTransactionController)
		}

		budgets := Balances.Group("/budgets")
		{
			budgets.GET("", mod.Budget.Ctl.ListBudgetController)
			budgets.POST("", requireMemberJWT(mod), ownerBudgetCreateMiddleware(mod), mod.Budget.Ctl.CreateBudgetController)
			budgets.GET("/:id", requireMemberJWT(mod), ownerBudgetReadMiddleware(mod), mod.Budget.Ctl.InfoBudgetController)
			budgets.PATCH("/:id", requireMemberJWT(mod), ownerBudgetUpdateMiddleware(mod), mod.Budget.Ctl.UpdateBudgetController)
			budgets.DELETE("/:id", requireMemberJWT(mod), ownerBudgetDeleteMiddleware(mod), mod.Budget.Ctl.DeleteBudgetController)
		}
	}
}

func apiPublic(r *gin.RouterGroup, mod *modules.Modules) {
	publics := r.Group("/public")
	{
		auths := publics.Group("/auth")
		{
			auths.POST("/login", mod.Member.Ctl.LoginMemberController)
			auths.POST("/refresh", mod.Member.Ctl.RefreshMemberTokenController)
			auths.POST("/register", mod.Member.Ctl.RegisterMemberController)
		}
	}
}
