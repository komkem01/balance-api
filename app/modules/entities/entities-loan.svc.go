package entities

import (
	"balance/app/modules/entities/ent"
	entitiesinf "balance/app/modules/entities/inf"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ entitiesinf.LoanEntity = (*Service)(nil)

func (s *Service) CreateLoan(ctx context.Context, memberID *string, name string, lender string, totalAmount float64, remainingBalance float64, monthlyPayment float64, interestRate float64, colorCode string, startDate *time.Time, endDate *time.Time) (*ent.LoanEntity, error) {
	var mid *uuid.UUID
	if memberID != nil {
		v := strings.TrimSpace(*memberID)
		if v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			mid = &id
		}
	}

	model := &ent.LoanEntity{
		ID:               uuid.New(),
		MemberID:         mid,
		Name:             strings.TrimSpace(name),
		Lender:           strings.TrimSpace(lender),
		TotalAmount:      totalAmount,
		RemainingBalance: remainingBalance,
		MonthlyPayment:   monthlyPayment,
		InterestRate:     interestRate,
		ColorCode:        strings.TrimSpace(colorCode),
		StartDate:        startDate,
		EndDate:          endDate,
	}

	_, err := s.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) GetLoanByID(ctx context.Context, id string) (*ent.LoanEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.LoanEntity{}
	if err := s.db.NewSelect().Model(model).Where("loan.id = ?", uid).Where("loan.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) UpdateLoan(ctx context.Context, id string, name *string, lender *string, totalAmount *float64, remainingBalance *float64, monthlyPayment *float64, interestRate *float64, colorCode *string, startDate *time.Time, endDate *time.Time) (*ent.LoanEntity, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	model := &ent.LoanEntity{}
	if err := s.db.NewSelect().Model(model).Where("loan.id = ?", uid).Where("loan.deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}

	if name != nil {
		model.Name = strings.TrimSpace(*name)
	}
	if lender != nil {
		model.Lender = strings.TrimSpace(*lender)
	}
	if totalAmount != nil {
		model.TotalAmount = *totalAmount
	}
	if remainingBalance != nil {
		model.RemainingBalance = *remainingBalance
	}
	if monthlyPayment != nil {
		model.MonthlyPayment = *monthlyPayment
	}
	if interestRate != nil {
		model.InterestRate = *interestRate
	}
	if colorCode != nil {
		model.ColorCode = strings.TrimSpace(*colorCode)
	}
	if startDate != nil {
		model.StartDate = startDate
	}
	if endDate != nil {
		model.EndDate = endDate
	}
	model.UpdatedAt = time.Now()

	_, err = s.db.NewUpdate().
		Model(model).
		WherePK().
		Column("name", "lender", "total_amount", "remaining_balance", "monthly_payment", "interest_rate", "color_code", "start_date", "end_date", "updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *Service) DeleteLoan(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	model := &ent.LoanEntity{ID: uid}
	_, err = s.db.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (s *Service) ListLoans(ctx context.Context, memberID *string) ([]*ent.LoanEntity, error) {
	items := make([]*ent.LoanEntity, 0)
	q := s.db.NewSelect().Model(&items).Where("loan.deleted_at IS NULL").Order("loan.created_at DESC")
	if memberID != nil {
		v := strings.TrimSpace(*memberID)
		if v != "" {
			mid, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			q = q.Where("loan.member_id = ?", mid)
		}
	}
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
