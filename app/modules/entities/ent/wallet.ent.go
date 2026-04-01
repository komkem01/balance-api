package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type WalletEntity struct {
	bun.BaseModel `bun:"table:wallets,alias:wallet"`

	ID        uuid.UUID  `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID  *uuid.UUID `bun:"type:uuid"`
	Name      string     `bun:"type:varchar"`
	Balance   float64    `bun:"type:numeric(18,2),notnull,default:0"`
	Currency  string     `bun:"type:varchar,notnull,default:'THB'"`
	ColorCode string     `bun:"type:varchar"`
	IsActive  bool       `bun:"type:boolean,notnull,default:true"`
	CreatedAt time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt time.Time  `bun:"type:timestamptz,notnull,default:now()"`
}
