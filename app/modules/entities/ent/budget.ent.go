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

	ID          uuid.UUID    `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID    *uuid.UUID   `bun:"type:uuid"`
	CategoryID  *uuid.UUID   `bun:"type:uuid"`
	Amount      float64      `bun:"type:numeric(18,2),notnull,default:0"`
	SpentAmount float64      `bun:"type:numeric(18,2),notnull,default:0"`
	Period      BudgetPeriod `bun:"type:budget_period,notnull"`
	StartDate   *time.Time   `bun:"type:date"`
	EndDate     *time.Time   `bun:"type:date"`
	CreatedAt   time.Time    `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt   time.Time    `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt   *time.Time   `bun:"type:timestamptz,soft_delete,nullzero"`
}
