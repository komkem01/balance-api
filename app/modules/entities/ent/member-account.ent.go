package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberAccountEntity struct {
	bun.BaseModel `bun:"table:member_accounts,alias:member_account"`

	ID        uuid.UUID  `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID  *uuid.UUID `bun:"type:uuid"`
	Username  string     `bun:"type:varchar"`
	Password  string     `bun:"type:text"`
	CreatedAt time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt time.Time  `bun:"type:timestamptz,soft_delete,nullzero"`
}
