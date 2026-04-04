package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberNotificationEntity = (*Service)(nil)

func (s *Service) CreateMemberNotification(ctx context.Context, memberID string, notificationType ent.MemberNotificationType, level ent.MemberNotificationLevel, title string, description string, dedupeKey *string) (*ent.MemberNotificationEntity, error) {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return nil, err
	}

	normalizedDedupe := ""
	if dedupeKey != nil {
		normalizedDedupe = strings.TrimSpace(*dedupeKey)
	}

	if normalizedDedupe != "" {
		exists := &ent.MemberNotificationEntity{}
		err := s.db.NewSelect().
			Model(exists).
			Where("member_notification.member_id = ?", mID).
			Where("member_notification.type = ?", notificationType).
			Where("member_notification.dedupe_key = ?", normalizedDedupe).
			Where("member_notification.deleted_at IS NULL").
			Scan(ctx)
		if err == nil {
			return exists, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	model := &ent.MemberNotificationEntity{
		ID:          uuid.New(),
		MemberID:    mID,
		Type:        notificationType,
		Level:       level,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		DedupeKey:   normalizedDedupe,
		IsRead:      false,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) ListMemberNotifications(ctx context.Context, memberID string, includeRead bool, limit int) ([]*ent.MemberNotificationEntity, error) {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	items := make([]*ent.MemberNotificationEntity, 0)
	q := s.db.NewSelect().
		Model(&items).
		Where("member_notification.member_id = ?", mID).
		Where("member_notification.deleted_at IS NULL").
		Order("member_notification.created_at DESC").
		Limit(limit)

	if !includeRead {
		q = q.Where("member_notification.is_read = FALSE")
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) GetLatestMemberNotificationByType(ctx context.Context, memberID string, notificationType ent.MemberNotificationType) (*ent.MemberNotificationEntity, error) {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return nil, err
	}

	item := &ent.MemberNotificationEntity{}
	if err := s.db.NewSelect().
		Model(item).
		Where("member_notification.member_id = ?", mID).
		Where("member_notification.type = ?", notificationType).
		Where("member_notification.deleted_at IS NULL").
		Order("member_notification.created_at DESC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) SetMemberNotificationRead(ctx context.Context, memberID string, notificationID string, isRead bool) (*ent.MemberNotificationEntity, error) {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return nil, err
	}

	nID, err := uuid.Parse(strings.TrimSpace(notificationID))
	if err != nil {
		return nil, err
	}

	item := &ent.MemberNotificationEntity{}
	if err := s.db.NewSelect().
		Model(item).
		Where("member_notification.id = ?", nID).
		Where("member_notification.member_id = ?", mID).
		Where("member_notification.deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	item.IsRead = isRead
	if isRead {
		now := time.Now().UTC()
		item.ReadAt = &now
	} else {
		item.ReadAt = nil
	}

	_, err = s.db.NewUpdate().
		Model(item).
		WherePK().
		Column("is_read", "read_at", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) SetAllMemberNotificationsRead(ctx context.Context, memberID string, isRead bool) error {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return err
	}

	q := s.db.NewUpdate().
		Model(&ent.MemberNotificationEntity{}).
		Set("is_read = ?", isRead).
		Set("updated_at = now()").
		Where("member_id = ?", mID).
		Where("deleted_at IS NULL")

	if isRead {
		q = q.Set("read_at = now()")
	} else {
		q = q.Set("read_at = NULL")
	}

	_, err = q.Exec(ctx)
	return err
}

func (s *Service) ClearMemberNotifications(ctx context.Context, memberID string) error {
	mID, err := uuid.Parse(strings.TrimSpace(memberID))
	if err != nil {
		return err
	}

	_, err = s.db.NewUpdate().
		Model(&ent.MemberNotificationEntity{}).
		Set("deleted_at = now()").
		Set("updated_at = now()").
		Where("member_id = ?", mID).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (s *Service) CreateMemberNotificationTx(ctx context.Context, tx bun.Tx, memberID uuid.UUID, notificationType ent.MemberNotificationType, level ent.MemberNotificationLevel, title string, description string, dedupeKey *string) (*ent.MemberNotificationEntity, error) {
	normalizedDedupe := ""
	if dedupeKey != nil {
		normalizedDedupe = strings.TrimSpace(*dedupeKey)
	}

	if normalizedDedupe != "" {
		exists := &ent.MemberNotificationEntity{}
		err := tx.NewSelect().
			Model(exists).
			Where("member_notification.member_id = ?", memberID).
			Where("member_notification.type = ?", notificationType).
			Where("member_notification.dedupe_key = ?", normalizedDedupe).
			Where("member_notification.deleted_at IS NULL").
			Scan(ctx)
		if err == nil {
			return exists, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	model := &ent.MemberNotificationEntity{
		ID:          uuid.New(),
		MemberID:    memberID,
		Type:        notificationType,
		Level:       level,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		DedupeKey:   normalizedDedupe,
		IsRead:      false,
	}

	if _, err := tx.NewInsert().Model(model).Exec(ctx); err != nil {
		return nil, err
	}

	return model, nil
}
