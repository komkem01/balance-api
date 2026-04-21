package loans

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateRequestService struct {
	MemberID         *string `json:"member_id"`
	Name             string  `json:"name"`
	Lender           string  `json:"lender"`
	TotalAmount      float64 `json:"total_amount"`
	RemainingBalance float64 `json:"remaining_balance"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	InterestRate     float64 `json:"interest_rate"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
}

type CreateResponseService struct {
	ID               uuid.UUID  `json:"id"`
	MemberID         *uuid.UUID `json:"member_id"`
	Name             string     `json:"name"`
	Lender           string     `json:"lender"`
	TotalAmount      float64    `json:"total_amount"`
	RemainingBalance float64    `json:"remaining_balance"`
	MonthlyPayment   float64    `json:"monthly_payment"`
	InterestRate     float64    `json:"interest_rate"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func parseLoanDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Service) CreateLoan(ctx context.Context, req *CreateRequestService) (*CreateResponseService, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrLoanNameRequired
	}

	if req.MemberID != nil {
		v := strings.TrimSpace(*req.MemberID)
		if v != "" {
			if _, err := uuid.Parse(v); err != nil {
				return nil, ErrLoanInvalidMemberID
			}
			if _, err := s.db.GetMemberByID(ctx, v); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, ErrLoanInvalidMemberID
				}
				return nil, err
			}
		}
	}

	startDate, err := parseLoanDate(req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := parseLoanDate(req.EndDate)
	if err != nil {
		return nil, err
	}

	item, err := s.db.CreateLoan(ctx, req.MemberID, req.Name, req.Lender, req.TotalAmount, req.RemainingBalance, req.MonthlyPayment, req.InterestRate, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return &CreateResponseService{
		ID:               item.ID,
		MemberID:         item.MemberID,
		Name:             item.Name,
		Lender:           item.Lender,
		TotalAmount:      item.TotalAmount,
		RemainingBalance: item.RemainingBalance,
		MonthlyPayment:   item.MonthlyPayment,
		InterestRate:     item.InterestRate,
		StartDate:        item.StartDate,
		EndDate:          item.EndDate,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}, nil
}
