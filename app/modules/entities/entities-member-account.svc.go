package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ entitiesinf.MemberAccountEntity = (*Service)(nil)

func parseOptionalAccountMemberID(value *string) (*uuid.UUID, error) {
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

func (s *Service) CreateMemberAccount(ctx context.Context, memberID *string, username string, password string) (*ent.MemberAccountEntity, error) {
	mid, err := parseOptionalAccountMemberID(memberID)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberAccountEntity{
		ID:       uuid.New(),
		MemberID: mid,
		Username: strings.TrimSpace(username),
		Password: password,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) GetMemberAccountByID(ctx context.Context, id string) (*ent.MemberAccountEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberAccountEntity{}
	if err := s.db.NewSelect().
		Model(model).
		Where("member_account.id = ?", uid).
		Scan(ctx); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) UpdateMemberAccount(ctx context.Context, id string, memberID *string, username *string, password *string) (*ent.MemberAccountEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	model := &ent.MemberAccountEntity{}
	if err := s.db.NewSelect().
		Model(model).
		Where("member_account.id = ?", uid).
		Scan(ctx); err != nil {
		return nil, err
	}

	if memberID != nil {
		parsed, err := parseOptionalAccountMemberID(memberID)
		if err != nil {
			return nil, err
		}
		model.MemberID = parsed
	}
	if username != nil {
		model.Username = strings.TrimSpace(*username)
	}
	if password != nil {
		model.Password = *password
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("member_id", "username", "password", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) DeleteMemberAccount(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = s.db.NewDelete().
		Model(&ent.MemberAccountEntity{}).
		Where("id = ?", uid).
		Exec(ctx)
	return err
}

func (s *Service) ListMemberAccounts(ctx context.Context) ([]*ent.MemberAccountEntity, error) {
	items := make([]*ent.MemberAccountEntity, 0)

	if err := s.db.NewSelect().
		Model(&items).
		Order("member_account.created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
