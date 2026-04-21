package categories

import "errors"

var (
	ErrCategoryNotFound         = errors.New("category not found")
	ErrCategoryInvalidID        = errors.New("invalid category ID")
	ErrCategoryInvalidMemberID  = errors.New("invalid category member ID")
	ErrCategoryNameRequired     = errors.New("category name is required")
	ErrCategoryTypeInvalid      = errors.New("category type is invalid")
	ErrCategoryPurposeInvalid   = errors.New("category purpose is invalid")
	ErrCategoryNoFieldsToUpdate = errors.New("no fields to update")
)
