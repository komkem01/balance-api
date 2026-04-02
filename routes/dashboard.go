package routes

import (
	"net/http"
	"sort"
	"time"

	"balance/app/modules"
	"balance/app/modules/entities/ent"
	"balance/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type dashboardSummaryResponse struct {
	Stats              dashboardStats             `json:"stats"`
	Wallets            []dashboardWalletItem      `json:"wallets"`
	RecentTransactions []dashboardTransactionItem `json:"recent_transactions"`
	ActiveBudgets      []dashboardBudgetItem      `json:"active_budgets"`
	GeneratedAt        time.Time                  `json:"generated_at"`
}

type dashboardStats struct {
	TotalBalance           float64 `json:"total_balance"`
	TotalIncome            float64 `json:"total_income"`
	TotalExpense           float64 `json:"total_expense"`
	NetAmount              float64 `json:"net_amount"`
	WalletCount            int     `json:"wallet_count"`
	ActiveWalletCount      int     `json:"active_wallet_count"`
	BudgetCount            int     `json:"budget_count"`
	RecentTransactionCount int     `json:"recent_transaction_count"`
}

type dashboardWalletItem struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
	IsActive bool      `json:"is_active"`
}

type dashboardTransactionItem struct {
	ID              uuid.UUID           `json:"id"`
	WalletID        *uuid.UUID          `json:"wallet_id"`
	CategoryID      *uuid.UUID          `json:"category_id"`
	Amount          float64             `json:"amount"`
	Type            ent.TransactionType `json:"type"`
	TransactionDate *time.Time          `json:"transaction_date"`
	Note            string              `json:"note"`
	CreatedAt       time.Time           `json:"created_at"`
}

type dashboardBudgetItem struct {
	ID         uuid.UUID        `json:"id"`
	CategoryID *uuid.UUID       `json:"category_id"`
	Amount     float64          `json:"amount"`
	Period     ent.BudgetPeriod `json:"period"`
	StartDate  *time.Time       `json:"start_date"`
	EndDate    *time.Time       `json:"end_date"`
}

func dashboardSummaryHandler(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberIDStr := resolveMemberIDFromContext(ctx)
		if memberIDStr == "" {
			_ = base.Unauthorized(ctx, "member-token-invalid", nil)
			return
		}
		memberID, err := uuid.Parse(memberIDStr)
		if err != nil {
			_ = base.Unauthorized(ctx, "member-token-invalid", nil)
			return
		}

		wallets, err := mod.ENT.Svc.ListWallets(ctx, nil)
		if err != nil {
			_ = base.InternalServerError(ctx, "dashboard-wallets-list-failed", nil)
			return
		}

		memberWallets := make([]*ent.WalletEntity, 0)
		walletIDs := make(map[uuid.UUID]struct{})
		totalBalance := 0.0
		activeWalletCount := 0
		for _, wallet := range wallets {
			if wallet.MemberID == nil || *wallet.MemberID != memberID {
				continue
			}
			memberWallets = append(memberWallets, wallet)
			walletIDs[wallet.ID] = struct{}{}
			totalBalance += wallet.Balance
			if wallet.IsActive {
				activeWalletCount++
			}
		}

		budgets, err := mod.ENT.Svc.ListBudgets(ctx, &memberIDStr, nil, nil)
		if err != nil {
			_ = base.InternalServerError(ctx, "dashboard-budgets-list-failed", nil)
			return
		}

		allTransactions := make([]*ent.TransactionEntity, 0)
		totalIncome := 0.0
		totalExpense := 0.0
		for _, wallet := range memberWallets {
			walletID := wallet.ID.String()
			items, txErr := mod.ENT.Svc.ListTransactions(ctx, nil, &walletID, nil, nil)
			if txErr != nil {
				_ = base.InternalServerError(ctx, "dashboard-transactions-list-failed", nil)
				return
			}
			for _, tx := range items {
				if tx.WalletID == nil {
					continue
				}
				if _, exists := walletIDs[*tx.WalletID]; !exists {
					continue
				}
				allTransactions = append(allTransactions, tx)
				switch tx.Type {
				case ent.TransactionTypeIncome:
					totalIncome += tx.Amount
				case ent.TransactionTypeExpense:
					totalExpense += tx.Amount
				}
			}
		}

		sort.Slice(allTransactions, func(i, j int) bool {
			return compareDashboardTransactionTime(allTransactions[i]).After(compareDashboardTransactionTime(allTransactions[j]))
		})

		recentLimit := 5
		if len(allTransactions) < recentLimit {
			recentLimit = len(allTransactions)
		}

		walletItems := make([]dashboardWalletItem, 0, len(memberWallets))
		for _, wallet := range memberWallets {
			walletItems = append(walletItems, dashboardWalletItem{
				ID:       wallet.ID,
				Name:     wallet.Name,
				Balance:  wallet.Balance,
				Currency: wallet.Currency,
				IsActive: wallet.IsActive,
			})
		}

		recentItems := make([]dashboardTransactionItem, 0, recentLimit)
		for _, tx := range allTransactions[:recentLimit] {
			recentItems = append(recentItems, dashboardTransactionItem{
				ID:              tx.ID,
				WalletID:        tx.WalletID,
				CategoryID:      tx.CategoryID,
				Amount:          tx.Amount,
				Type:            tx.Type,
				TransactionDate: tx.TransactionDate,
				Note:            tx.Note,
				CreatedAt:       tx.CreatedAt,
			})
		}

		budgetItems := make([]dashboardBudgetItem, 0, len(budgets))
		for _, budget := range budgets {
			budgetItems = append(budgetItems, dashboardBudgetItem{
				ID:         budget.ID,
				CategoryID: budget.CategoryID,
				Amount:     budget.Amount,
				Period:     budget.Period,
				StartDate:  budget.StartDate,
				EndDate:    budget.EndDate,
			})
		}

		res := dashboardSummaryResponse{
			Stats: dashboardStats{
				TotalBalance:           totalBalance,
				TotalIncome:            totalIncome,
				TotalExpense:           totalExpense,
				NetAmount:              totalIncome - totalExpense,
				WalletCount:            len(memberWallets),
				ActiveWalletCount:      activeWalletCount,
				BudgetCount:            len(budgetItems),
				RecentTransactionCount: len(recentItems),
			},
			Wallets:            walletItems,
			RecentTransactions: recentItems,
			ActiveBudgets:      budgetItems,
			GeneratedAt:        time.Now().UTC(),
		}

		ctx.Header("Cache-Control", "no-store")
		_ = base.RawJSON(ctx, http.StatusOK, res)
	}
}

func compareDashboardTransactionTime(tx *ent.TransactionEntity) time.Time {
	if tx.TransactionDate != nil {
		return *tx.TransactionDate
	}
	return tx.CreatedAt
}
