package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type LoanEntity struct {
	bun.BaseModel `bun:"table:loans,alias:loan"`

	ID               uuid.UUID  `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID         *uuid.UUID `bun:"type:uuid"`
	Name             string     `bun:"type:varchar,notnull"`
	Lender           string     `bun:"type:varchar"`
	TotalAmount      float64    `bun:"type:numeric,notnull,default:0"`
	RemainingBalance float64    `bun:"type:numeric,notnull,default:0"`
	MonthlyPayment   float64    `bun:"type:numeric,notnull,default:0"`
	InterestRate     float64    `bun:"type:numeric,notnull,default:0"`
	ColorCode        string     `bun:"type:varchar,notnull,default:'#6366f1'"`
	StartDate        *time.Time `bun:"type:date"`
	EndDate          *time.Time `bun:"type:date"`
	CreatedAt        time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt        time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt        *time.Time `bun:"type:timestamptz,soft_delete,nullzero"`
}
