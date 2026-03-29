package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PrefixEntity struct {
	bun.BaseModel `bun:"table:prefixes,alias:prefix"`

	ID        uuid.UUID `bun:"type:uuid,default:gen_random_uuid(),pk"`
	GenderID  uuid.UUID `bun:"type:uuid,notnull"`
	Name      string    `bun:"type:text,notnull"`
	IsActive  bool      `bun:"type:boolean,notnull,default:true"`
	CreatedAt time.Time `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt time.Time `bun:"type:timestamptz,notnull,default:now()"`
}
