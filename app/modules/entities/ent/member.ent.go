package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberEntity struct {
	bun.BaseModel `bun:"table:members,alias:member"`

	ID          uuid.UUID  `bun:"type:uuid,default:gen_random_uuid(),pk"`
	GenderID    *uuid.UUID `bun:"type:uuid"`
	PrefixID    *uuid.UUID `bun:"type:uuid"`
	FirstName   string     `bun:"type:varchar"`
	LastName    string     `bun:"type:varchar"`
	DisplayName string     `bun:"type:varchar"`
	Phone       string     `bun:"type:varchar"`
	CreatedAt   time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt   time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	LastLogin   *time.Time `bun:"type:timestamptz"`
	DeletedAt   time.Time  `bun:"type:timestamptz,soft_delete,nullzero"`
}
