package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GoalEntrySourceType string

const (
	GoalEntrySourceTypeTransaction GoalEntrySourceType = "transaction"
	GoalEntrySourceTypeLoan        GoalEntrySourceType = "loan"
	GoalEntrySourceTypeSystem      GoalEntrySourceType = "system"
)

type GoalEntryEntity struct {
	bun.BaseModel `bun:"table:goal_entries,alias:goal_entry"`

	ID           uuid.UUID           `bun:"type:uuid,default:gen_random_uuid(),pk"`
	GoalID       uuid.UUID           `bun:"type:uuid,notnull"`
	MemberID     *uuid.UUID          `bun:"type:uuid"`
	SourceType   GoalEntrySourceType `bun:"type:varchar,notnull"`
	SourceID     *uuid.UUID          `bun:"type:uuid"`
	AmountBefore float64             `bun:"type:numeric,notnull,default:0"`
	AmountAfter  float64             `bun:"type:numeric,notnull,default:0"`
	AmountDelta  float64             `bun:"type:numeric,notnull,default:0"`
	Note         string              `bun:"type:varchar,notnull,default:''"`
	CreatedAt    time.Time           `bun:"type:timestamptz,notnull,default:now()"`
}
