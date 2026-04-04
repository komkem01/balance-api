package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberNotificationType string

const (
	MemberNotificationTypeBudget   MemberNotificationType = "budget"
	MemberNotificationTypeSecurity MemberNotificationType = "security"
	MemberNotificationTypeWeekly   MemberNotificationType = "weekly"
)

type MemberNotificationLevel string

const (
	MemberNotificationLevelInfo     MemberNotificationLevel = "info"
	MemberNotificationLevelWarning  MemberNotificationLevel = "warning"
	MemberNotificationLevelCritical MemberNotificationLevel = "critical"
)

type MemberNotificationEntity struct {
	bun.BaseModel `bun:"table:member_notifications,alias:member_notification"`

	ID          uuid.UUID               `bun:"type:uuid,default:gen_random_uuid(),pk"`
	MemberID    uuid.UUID               `bun:"type:uuid,notnull"`
	Type        MemberNotificationType  `bun:"type:varchar(20),notnull"`
	Level       MemberNotificationLevel `bun:"type:varchar(20),notnull"`
	Title       string                  `bun:"type:varchar(255),notnull"`
	Description string                  `bun:"type:text,notnull"`
	DedupeKey   string                  `bun:"type:varchar(255)"`
	IsRead      bool                    `bun:"type:boolean,notnull,default:false"`
	ReadAt      *time.Time              `bun:"type:timestamptz"`
	CreatedAt   time.Time               `bun:"type:timestamptz,notnull,default:now()"`
	UpdatedAt   time.Time               `bun:"type:timestamptz,notnull,default:now()"`
	DeletedAt   *time.Time              `bun:"type:timestamptz,soft_delete,nullzero"`
}
