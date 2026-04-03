package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberEntity = (*Service)(nil)

func parseOptionalUUIDString(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}

	id, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *Service) CreateMember(ctx context.Context, genderID *string, prefixID *string, firstName string, lastName string, displayName string, phone string) (*ent.MemberEntity, error) {
	gid, err := parseOptionalUUIDString(genderID)
	if err != nil {
		return nil, err
	}

	pid, err := parseOptionalUUIDString(prefixID)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberEntity{
		ID:                uuid.New(),
		GenderID:          gid,
		PrefixID:          pid,
		FirstName:         strings.TrimSpace(firstName),
		LastName:          strings.TrimSpace(lastName),
		DisplayName:       strings.TrimSpace(displayName),
		Phone:             strings.TrimSpace(phone),
		PreferredCurrency: "THB",
		PreferredLanguage: "EN",
		NotifyBudget:      true,
		NotifySecurity:    true,
		NotifyWeekly:      false,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) CreateMemberWithAccount(ctx context.Context, genderID *string, prefixID *string, firstName string, lastName string, displayName string, phone string, username string, password string) (*ent.MemberEntity, error) {
	gid, err := parseOptionalUUIDString(genderID)
	if err != nil {
		return nil, err
	}

	pid, err := parseOptionalUUIDString(prefixID)
	if err != nil {
		return nil, err
	}

	member := &ent.MemberEntity{
		ID:                uuid.New(),
		GenderID:          gid,
		PrefixID:          pid,
		FirstName:         strings.TrimSpace(firstName),
		LastName:          strings.TrimSpace(lastName),
		DisplayName:       strings.TrimSpace(displayName),
		Phone:             strings.TrimSpace(phone),
		PreferredCurrency: "THB",
		PreferredLanguage: "EN",
		NotifyBudget:      true,
		NotifySecurity:    true,
		NotifyWeekly:      false,
	}

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return err
		}

		memberID := member.ID
		account := &ent.MemberAccountEntity{
			ID:       uuid.New(),
			MemberID: &memberID,
			Username: strings.TrimSpace(username),
			Password: password,
		}

		if _, err := tx.NewInsert().Model(account).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *Service) GetMemberByID(ctx context.Context, id string) (*ent.MemberEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberEntity{}
	if err := s.db.NewSelect().
		Model(model).
		Where("member.id = ?", uid).
		Scan(ctx); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) UpdateMember(ctx context.Context, id string, genderID *string, prefixID *string, firstName *string, lastName *string, displayName *string, phone *string, lastLogin *time.Time, profileImageURL *string) (*ent.MemberEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberEntity{}
	if err := s.db.NewSelect().
		Model(model).
		Where("member.id = ?", uid).
		Scan(ctx); err != nil {
		return nil, err
	}

	if genderID != nil {
		parsed, err := parseOptionalUUIDString(genderID)
		if err != nil {
			return nil, err
		}
		model.GenderID = parsed
	}

	if prefixID != nil {
		parsed, err := parseOptionalUUIDString(prefixID)
		if err != nil {
			return nil, err
		}
		model.PrefixID = parsed
	}

	if firstName != nil {
		model.FirstName = strings.TrimSpace(*firstName)
	}
	if lastName != nil {
		model.LastName = strings.TrimSpace(*lastName)
	}
	if displayName != nil {
		model.DisplayName = strings.TrimSpace(*displayName)
	}
	if phone != nil {
		model.Phone = strings.TrimSpace(*phone)
	}
	if lastLogin != nil {
		model.LastLogin = lastLogin
	}
	if profileImageURL != nil {
		model.ProfileImageURL = strings.TrimSpace(*profileImageURL)
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("gender_id", "prefix_id", "first_name", "last_name", "display_name", "phone", "last_login", "profile_image_url", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) UpdateMemberSettings(ctx context.Context, id string, preferredCurrency *string, preferredLanguage *string, notifyBudget *bool, notifySecurity *bool, notifyWeekly *bool) (*ent.MemberEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberEntity{}
	if err := s.db.NewSelect().
		Model(model).
		Where("member.id = ?", uid).
		Scan(ctx); err != nil {
		return nil, err
	}

	if preferredCurrency != nil {
		model.PreferredCurrency = strings.TrimSpace(*preferredCurrency)
	}
	if preferredLanguage != nil {
		model.PreferredLanguage = strings.TrimSpace(*preferredLanguage)
	}
	if notifyBudget != nil {
		model.NotifyBudget = *notifyBudget
	}
	if notifySecurity != nil {
		model.NotifySecurity = *notifySecurity
	}
	if notifyWeekly != nil {
		model.NotifyWeekly = *notifyWeekly
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("preferred_currency", "preferred_language", "notify_budget", "notify_security", "notify_weekly", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) DeleteMember(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = s.db.NewDelete().
		Model(&ent.MemberEntity{}).
		Where("id = ?", uid).
		Exec(ctx)
	return err
}

func (s *Service) DeleteMemberWithAccounts(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().
			Model(&ent.MemberAccountEntity{}).
			Where("member_id = ?", uid).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewDelete().
			Model(&ent.MemberEntity{}).
			Where("id = ?", uid).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ListMembers(ctx context.Context) ([]*ent.MemberEntity, error) {
	items := make([]*ent.MemberEntity, 0)

	if err := s.db.NewSelect().
		Model(&items).
		Order("member.created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
