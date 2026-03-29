package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BudgetPeriod string

const (
	BudgetPeriodDaily   BudgetPeriod = "daily"
	BudgetPeriodWeekly  BudgetPeriod = "weekly"
	BudgetPeriodMonthly BudgetPeriod = "monthly"
)

type BudgetEntity struct {
	bun.BaseModel `bun:"table:budgets,alias:budget"`

	ID         uuid.UUID    `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID   *uuid.UUID   `bun:"type:uuid"`
	CategoryID *uuid.UUID   `bun:"type:uuid"`
	Amount     float64      `bun:"type:numeric,notnull,default:0"`
	Period     BudgetPeriod `bun:"type:budget_period,notnull"`
	StartDate  *time.Time   `bun:"type:date"`
	EndDate    *time.Time   `bun:"type:date"`
	CreatedAt  time.Time    `bun:"type:timestamptz,notnull,default:now()"`
}
