package loans

import (
	"errors"
)

var (
	ErrLoanInvalidID        = errors.New("invalid loan ID")
	ErrLoanNameRequired     = errors.New("loan name is required")
	ErrLoanNotFound         = errors.New("loan not found")
	ErrLoanInvalidMemberID  = errors.New("invalid loan member ID")
	ErrLoanNoFieldsToUpdate = errors.New("no fields to update")
	ErrLoanStartDateInvalid = errors.New("invalid loan start date")
	ErrLoanEndDateInvalid   = errors.New("invalid loan end date")
)
