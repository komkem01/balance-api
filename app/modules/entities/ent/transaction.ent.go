package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type TransactionEntity struct {
	bun.BaseModel `bun:"table:transactions,alias:transaction"`

	ID              uuid.UUID       `bun:"type:uuid,default:gen_random_uuid(),pk"`
	WalletID        *uuid.UUID      `bun:"type:uuid"`
	CategoryID      *uuid.UUID      `bun:"type:uuid"`
	Amount          float64         `bun:"type:numeric(18,2),notnull,default:0"`
	Type            TransactionType `bun:"type:transaction_type,notnull"`
	TransactionDate *time.Time      `bun:"type:date"`
	Note            string          `bun:"type:text"`
	ImageURL        string          `bun:"type:varchar"`
	CreatedAt       time.Time       `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt       time.Time       `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt       *time.Time      `bun:"type:timestamptz,soft_delete,nullzero"`
}
