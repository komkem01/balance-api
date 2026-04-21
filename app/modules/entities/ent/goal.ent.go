package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GoalType string

const (
	GoalTypeSavings    GoalType = "savings"
	GoalTypeDebtPayoff GoalType = "debt_payoff"
)

type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusPaused    GoalStatus = "paused"
	GoalStatusArchived  GoalStatus = "archived"
)

type GoalTrackingSourceType string

const (
	GoalTrackingSourceWallet     GoalTrackingSourceType = "wallet"
	GoalTrackingSourceAllWallets GoalTrackingSourceType = "all_wallets"
	GoalTrackingSourceLoan       GoalTrackingSourceType = "loan"
)

type GoalEntity struct {
	bun.BaseModel `bun:"table:goals,alias:goal"`

	ID                 uuid.UUID               `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID           *uuid.UUID              `bun:"type:uuid"`
	Name               string                  `bun:"type:varchar,notnull"`
	Type               GoalType                `bun:"type:varchar,notnull"`
	TargetAmount       float64                 `bun:"type:numeric,notnull,default:0"`
	StartAmount        float64                 `bun:"type:numeric,notnull,default:0"`
	CurrentAmount      float64                 `bun:"type:numeric,notnull,default:0"`
	StartDate          *time.Time              `bun:"type:date"`
	TargetDate         *time.Time              `bun:"type:date"`
	Status             GoalStatus              `bun:"type:varchar,notnull,default:'active'"`
	AutoTracking       bool                    `bun:"type:boolean,notnull,default:true"`
	TrackingSourceType *GoalTrackingSourceType `bun:"type:varchar"`
	TrackingSourceID   *uuid.UUID              `bun:"type:uuid"`
	DepositWalletID    *uuid.UUID              `bun:"type:uuid"`
	CreatedAt          time.Time               `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt          time.Time               `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt          *time.Time              `bun:"type:timestamptz,soft_delete,nullzero"`
}
