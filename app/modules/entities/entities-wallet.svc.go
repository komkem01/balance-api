package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"

	"github.com/google/uuid"
)

var _ entitiesinf.WalletEntity = (*Service)(nil)

func parseWalletMemberID(value *string) (*uuid.UUID, error) {
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

func (s *Service) CreateWallet(ctx context.Context, memberID *string, name string, balance float64, currency string, colorCode string, isActive bool) (*ent.WalletEntity, error) {
	mid, err := parseWalletMemberID(memberID)
	if err != nil {
		return nil, err
	}

	model := &ent.WalletEntity{
		ID:        uuid.New(),
		MemberID:  mid,
		Name:      strings.TrimSpace(name),
		Balance:   balance,
		Currency:  strings.TrimSpace(currency),
		ColorCode: strings.TrimSpace(colorCode),
		IsActive:  isActive,
	}

	_, err = s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) GetWalletByID(ctx context.Context, id string) (*ent.WalletEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.WalletEntity{}
	if err := s.db.NewSelect().Model(model).Where("wallet.id = ?", uid).Where("wallet.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateWallet(ctx context.Context, id string, memberID *string, name *string, balance *float64, currency *string, colorCode *string, isActive *bool) (*ent.WalletEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.WalletEntity{}
	if err := s.db.NewSelect().Model(model).Where("wallet.id = ?", uid).Where("wallet.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}

	if memberID != nil {
		mid, err := parseWalletMemberID(memberID)
		if err != nil {
			return nil, err
		}
		model.MemberID = mid
	}
	if name != nil {
		model.Name = strings.TrimSpace(*name)
	}
	if balance != nil {
		model.Balance = *balance
	}
	if currency != nil {
		model.Currency = strings.TrimSpace(*currency)
	}
	if colorCode != nil {
		model.ColorCode = strings.TrimSpace(*colorCode)
	}
	if isActive != nil {
		model.IsActive = *isActive
	}

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("member_id", "name", "balance", "currency", "color_code", "is_active", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) DeleteWallet(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = s.db.NewUpdate().
		Model(&ent.WalletEntity{}).
		Set("deleted_at = now()").
		Set("updated_at = now()").
		Where("id = ?", uid).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (s *Service) ListWallets(ctx context.Context, isActive *bool) ([]*ent.WalletEntity, error) {
	items := make([]*ent.WalletEntity, 0)
	q := s.db.NewSelect().Model(&items).Where("wallet.deleted_at IS NULL").Order("wallet.created_at DESC")
	if isActive != nil {
		q = q.Where("wallet.is_active = ?", *isActive)
	}
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
