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
	ColorCode        *string `json:"color_code"`
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
	ColorCode        string     `json:"color_code"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

const defaultLoanColorCode = "#6366f1"

func normalizeLoanColorCode(value *string) string {
	if value == nil {
		return defaultLoanColorCode
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return defaultLoanColorCode
	}
	return v
}

func parseLoanDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}

	if t, err := time.Parse("2006-01-02", v); err == nil {
		return &t, nil
	} else {
		parseErr := err

		layouts := []string{time.RFC3339, time.RFC3339Nano}
		for _, layout := range layouts {
			t, err := time.Parse(layout, v)
			if err == nil {
				onlyDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
				return &onlyDate, nil
			}
			parseErr = err
		}

		return nil, parseErr
	}
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
		return nil, ErrLoanStartDateInvalid
	}
	endDate, err := parseLoanDate(req.EndDate)
	if err != nil {
		return nil, ErrLoanEndDateInvalid
	}

	item, err := s.db.CreateLoan(ctx, req.MemberID, req.Name, req.Lender, req.TotalAmount, req.RemainingBalance, req.MonthlyPayment, req.InterestRate, normalizeLoanColorCode(req.ColorCode), startDate, endDate)
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
		ColorCode:        item.ColorCode,
		StartDate:        item.StartDate,
		EndDate:          item.EndDate,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}, nil
}
