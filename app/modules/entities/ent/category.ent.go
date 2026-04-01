package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type CategoryEntity struct {
	bun.BaseModel `bun:"table:categories,alias:category"`

	ID        uuid.UUID    `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID  *uuid.UUID   `bun:"type:uuid"`
	Name      string       `bun:"type:varchar"`
	Type      CategoryType `bun:"type:category_type,notnull"`
	IconName  string       `bun:"type:varchar"`
	ColorCode string       `bun:"type:varchar"`
	CreatedAt time.Time    `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt time.Time    `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt *time.Time   `bun:"type:timestamptz,soft_delete,nullzero"`
}
