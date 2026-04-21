package loans

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateRequestService struct {
	ID               string
	Name             *string  `json:"name"`
	Lender           *string  `json:"lender"`
	TotalAmount      *float64 `json:"total_amount"`
	RemainingBalance *float64 `json:"remaining_balance"`
	MonthlyPayment   *float64 `json:"monthly_payment"`
	InterestRate     *float64 `json:"interest_rate"`
	StartDate        *string  `json:"start_date"`
	EndDate          *string  `json:"end_date"`
}

type UpdateResponseService struct {
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

func (s *Service) UpdateLoan(ctx context.Context, req *UpdateRequestService) (*UpdateResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrLoanInvalidID
	}
	if req.Name == nil && req.Lender == nil && req.TotalAmount == nil && req.RemainingBalance == nil && req.MonthlyPayment == nil && req.InterestRate == nil && req.StartDate == nil && req.EndDate == nil {
		return nil, ErrLoanNoFieldsToUpdate
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return nil, ErrLoanNameRequired
	}

	var startDate *time.Time
	if req.StartDate != nil {
		d, err := parseLoanDate(req.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = d
	}
	var endDate *time.Time
	if req.EndDate != nil {
		d, err := parseLoanDate(req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = d
	}

	item, err := s.db.UpdateLoan(ctx, req.ID, req.Name, req.Lender, req.TotalAmount, req.RemainingBalance, req.MonthlyPayment, req.InterestRate, startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}
	return &UpdateResponseService{
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
