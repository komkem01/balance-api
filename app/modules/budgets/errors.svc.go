package budgets

import "errors"

var (
	ErrBudgetNotFound          = errors.New("budget not found")
	ErrBudgetInvalidID         = errors.New("invalid budget ID")
	ErrBudgetInvalidMemberID   = errors.New("invalid budget member ID")
	ErrBudgetInvalidCategoryID = errors.New("invalid budget category ID")
	ErrBudgetPeriodInvalid     = errors.New("budget period is invalid")
	ErrBudgetNoFieldsToUpdate  = errors.New("no fields to update")
	ErrBudgetDateInvalid       = errors.New("budget date is invalid")
)
