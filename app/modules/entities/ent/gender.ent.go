package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GenderEntity struct {
	bun.BaseModel `bun:"table:genders,alias:gender"`

	ID        uuid.UUID `bun:"type:uuid,default:gen_random_uuid(),pk"`
	Name      string    `bun:"type:text,notnull"`
	IsActive  bool      `bun:"type:boolean,notnull,default:true"`
	CreatedAt time.Time `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt time.Time `bun:"type:timestamptz,notnull,default:now()"`
}
