package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberEntity struct {
	bun.BaseModel `bun:"table:members,alias:member"`

	ID                uuid.UUID  `bun:"type:uuid,default:gen_random_uuid(),pk"`
	GenderID          *uuid.UUID `bun:"type:uuid"`
	PrefixID          *uuid.UUID `bun:"type:uuid"`
	FirstName         string     `bun:"type:varchar"`
	LastName          string     `bun:"type:varchar"`
	DisplayName       string     `bun:"type:varchar"`
	Phone             string     `bun:"type:varchar"`
	ProfileImageURL   string     `bun:"type:varchar,notnull,default:''"`
	PreferredCurrency string     `bun:"type:varchar,notnull,default:'THB'"`
	PreferredLanguage string     `bun:"type:varchar,notnull,default:'EN'"`
	NotifyBudget      bool       `bun:"type:boolean,notnull,default:true"`
	NotifySecurity    bool       `bun:"type:boolean,notnull,default:true"`
	NotifyWeekly      bool       `bun:"type:boolean,notnull,default:false"`
	CreatedAt         time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt         time.Time  `bun:"type:timestamptz,notnull,default:now()"`
	LastLogin         *time.Time `bun:"type:timestamptz"`
	DeletedAt         time.Time  `bun:"type:timestamptz,soft_delete,nullzero"`
}
