package loans

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type InfoRequestService struct {
	ID string
}

type InfoResponseService struct {
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

func (s *Service) InfoLoan(ctx context.Context, req *InfoRequestService) (*InfoResponseService, error) {
	if _, err := uuid.Parse(req.ID); err != nil {
		return nil, ErrLoanInvalidID
	}
	item, err := s.db.GetLoanByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}
	return &InfoResponseService{
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
