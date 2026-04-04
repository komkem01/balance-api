package members

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

type MeNotificationItemService struct {
	ID          uuid.UUID                   `json:"id"`
	Type        ent.MemberNotificationType  `json:"type"`
	Level       ent.MemberNotificationLevel `json:"level"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	IsRead      bool                        `json:"is_read"`
	ReadAt      *time.Time                  `json:"read_at"`
	CreatedAt   time.Time                   `json:"created_at"`
}

type ListMeNotificationsRequestService struct {
	MemberID    string `json:"member_id"`
	IncludeRead bool   `json:"include_read"`
	Limit       int    `json:"limit"`
}

type SetMeNotificationReadRequestService struct {
	MemberID       string `json:"member_id"`
	NotificationID string `json:"notification_id"`
	IsRead         bool   `json:"is_read"`
}

func (s *Service) ensureWeeklyNotification(ctx context.Context, memberID string, preferredLanguage string) error {
	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		return err
	}
	if !member.NotifyWeekly {
		return nil
	}

	now := time.Now().UTC()
	endDate := now.Format("2006-01-02")
	startWindow := now.AddDate(0, 0, -6)
	startDate := time.Date(startWindow.Year(), startWindow.Month(), startWindow.Day(), 0, 0, 0, 0, time.UTC)
	endDateTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)

	transactions, err := s.db.ListTransactions(ctx, &memberID, nil, nil, nil)
	if err != nil {
		return err
	}

	var incomeTotal float64
	var expenseTotal float64
	for _, item := range transactions {
		referenceDate := item.CreatedAt.UTC()
		if item.TransactionDate != nil {
			referenceDate = item.TransactionDate.UTC()
		}

		if referenceDate.Before(startDate) || referenceDate.After(endDateTime) {
			continue
		}

		if strings.HasPrefix(item.Note, "__transfer__|") || item.Note == "Wallet transfer" {
			continue
		}

		if item.Type == ent.TransactionTypeIncome {
			incomeTotal += item.Amount
		} else {
			expenseTotal += item.Amount
		}
	}

	dedupeKey := fmt.Sprintf("weekly:%s", endDate)
	title := "Weekly financial summary"
	description := fmt.Sprintf("Income %.2f | Expense %.2f (last 7 days)", incomeTotal, expenseTotal)
	if strings.EqualFold(preferredLanguage, "TH") {
		title = "สรุปการเงิน 7 วันล่าสุด"
		description = fmt.Sprintf("รายรับ %.2f | รายจ่าย %.2f (ย้อนหลัง 7 วัน)", incomeTotal, expenseTotal)
	}

	_, err = s.db.CreateMemberNotification(
		ctx,
		memberID,
		ent.MemberNotificationTypeWeekly,
		ent.MemberNotificationLevelInfo,
		title,
		description,
		&dedupeKey,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListMeNotifications(ctx context.Context, req *ListMeNotificationsRequestService) ([]*MeNotificationItemService, error) {
	memberID := strings.TrimSpace(req.MemberID)
	if memberID == "" {
		return nil, ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return nil, ErrMemberUnauthorized
	}

	member, err := s.db.GetMemberByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	if err := s.ensureWeeklyNotification(ctx, memberID, member.PreferredLanguage); err != nil {
		return nil, err
	}

	items, err := s.db.ListMemberNotifications(ctx, memberID, req.IncludeRead, req.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]*MeNotificationItemService, 0, len(items))
	for _, item := range items {
		result = append(result, &MeNotificationItemService{
			ID:          item.ID,
			Type:        item.Type,
			Level:       item.Level,
			Title:       item.Title,
			Description: item.Description,
			IsRead:      item.IsRead,
			ReadAt:      item.ReadAt,
			CreatedAt:   item.CreatedAt,
		})
	}

	return result, nil
}

func (s *Service) SetMeNotificationRead(ctx context.Context, req *SetMeNotificationReadRequestService) (*MeNotificationItemService, error) {
	memberID := strings.TrimSpace(req.MemberID)
	if memberID == "" {
		return nil, ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return nil, ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(strings.TrimSpace(req.NotificationID)); err != nil {
		return nil, ErrMemberNotificationInvalidID
	}

	item, err := s.db.SetMemberNotificationRead(ctx, memberID, req.NotificationID, req.IsRead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotificationNotFound
		}
		return nil, err
	}

	return &MeNotificationItemService{
		ID:          item.ID,
		Type:        item.Type,
		Level:       item.Level,
		Title:       item.Title,
		Description: item.Description,
		IsRead:      item.IsRead,
		ReadAt:      item.ReadAt,
		CreatedAt:   item.CreatedAt,
	}, nil
}

func (s *Service) MarkAllMeNotificationsRead(ctx context.Context, memberID string) error {
	memberID = strings.TrimSpace(memberID)
	if memberID == "" {
		return ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return ErrMemberUnauthorized
	}

	if err := s.db.SetAllMemberNotificationsRead(ctx, memberID, true); err != nil {
		return err
	}

	return nil
}

func (s *Service) ClearMeNotifications(ctx context.Context, memberID string) error {
	memberID = strings.TrimSpace(memberID)
	if memberID == "" {
		return ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return ErrMemberUnauthorized
	}

	if err := s.db.ClearMemberNotifications(ctx, memberID); err != nil {
		return err
	}

	return nil
}
