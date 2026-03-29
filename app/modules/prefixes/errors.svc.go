package prefixes

import "errors"

var (
	ErrPrefixNotFound         = errors.New("prefix not found")
	ErrPrefixNameRequired     = errors.New("prefix name is required")
	ErrPrefixInvalidID        = errors.New("invalid prefix ID")
	ErrPrefixInvalidGenderID  = errors.New("invalid prefix gender ID")
	ErrPrefixNoFieldsToUpdate = errors.New("no fields to update")
	ErrPrefixAlreadyExists    = errors.New("prefix with the same name already exists")
)
