package ent

import "time"

type TransactionMonthlySummaryEntity struct {
	Month            time.Time `bun:"month"`
	IncomeTotal      float64   `bun:"income_total"`
	ExpenseTotal     float64   `bun:"expense_total"`
	TransactionCount int64     `bun:"transaction_count"`
}
