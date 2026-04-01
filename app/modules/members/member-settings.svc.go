package members

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type MeSettingsResponseService struct {
	PreferredCurrency string `json:"preferred_currency"`
	PreferredLanguage string `json:"preferred_language"`
	NotifyBudget      bool   `json:"notify_budget"`
	NotifySecurity    bool   `json:"notify_security"`
	NotifyWeekly      bool   `json:"notify_weekly"`
}

type MeSettingsUpdateRequestService struct {
	MemberID          string  `json:"member_id"`
	PreferredCurrency *string `json:"preferred_currency"`
	PreferredLanguage *string `json:"preferred_language"`
	NotifyBudget      *bool   `json:"notify_budget"`
	NotifySecurity    *bool   `json:"notify_security"`
	NotifyWeekly      *bool   `json:"notify_weekly"`
}

func normalizeCurrency(value string) (string, bool) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "THB", "USD", "EUR":
		return strings.ToUpper(strings.TrimSpace(value)), true
	default:
		return "", false
	}
}

func normalizeLanguage(value string) (string, bool) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "EN", "TH":
		return strings.ToUpper(strings.TrimSpace(value)), true
	default:
		return "", false
	}
}

func (s *Service) InfoMeSettings(ctx context.Context, memberID string) (*MeSettingsResponseService, error) {
	memberID = strings.TrimSpace(memberID)
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

	currency := strings.ToUpper(strings.TrimSpace(member.PreferredCurrency))
	if currency == "" {
		currency = "THB"
	}
	language := strings.ToUpper(strings.TrimSpace(member.PreferredLanguage))
	if language == "" {
		language = "EN"
	}

	return &MeSettingsResponseService{
		PreferredCurrency: currency,
		PreferredLanguage: language,
		NotifyBudget:      member.NotifyBudget,
		NotifySecurity:    member.NotifySecurity,
		NotifyWeekly:      member.NotifyWeekly,
	}, nil
}

func (s *Service) UpdateMeSettings(ctx context.Context, req *MeSettingsUpdateRequestService) (*MeSettingsResponseService, error) {
	memberID := strings.TrimSpace(req.MemberID)
	if memberID == "" {
		return nil, ErrMemberUnauthorized
	}
	if _, err := uuid.Parse(memberID); err != nil {
		return nil, ErrMemberUnauthorized
	}

	if req.PreferredCurrency == nil && req.PreferredLanguage == nil && req.NotifyBudget == nil && req.NotifySecurity == nil && req.NotifyWeekly == nil {
		return nil, ErrMemberNoSettingsToUpdate
	}

	var normalizedCurrency *string
	if req.PreferredCurrency != nil {
		v, ok := normalizeCurrency(*req.PreferredCurrency)
		if !ok {
			return nil, ErrMemberInvalidCurrency
		}
		normalizedCurrency = &v
	}

	var normalizedLanguage *string
	if req.PreferredLanguage != nil {
		v, ok := normalizeLanguage(*req.PreferredLanguage)
		if !ok {
			return nil, ErrMemberInvalidLanguage
		}
		normalizedLanguage = &v
	}

	member, err := s.db.UpdateMemberSettings(
		ctx,
		memberID,
		normalizedCurrency,
		normalizedLanguage,
		req.NotifyBudget,
		req.NotifySecurity,
		req.NotifyWeekly,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(member.PreferredCurrency))
	if currency == "" {
		currency = "THB"
	}
	language := strings.ToUpper(strings.TrimSpace(member.PreferredLanguage))
	if language == "" {
		language = "EN"
	}

	return &MeSettingsResponseService{
		PreferredCurrency: currency,
		PreferredLanguage: language,
		NotifyBudget:      member.NotifyBudget,
		NotifySecurity:    member.NotifySecurity,
		NotifyWeekly:      member.NotifyWeekly,
	}, nil
}
