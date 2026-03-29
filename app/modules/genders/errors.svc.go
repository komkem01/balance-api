package genders

import "errors"

var (
	ErrGenderNotFound         = errors.New("gender not found")
	ErrGenderNameRequired     = errors.New("gender name is required")
	ErrGenderInvalidID        = errors.New("invalid gender ID")
	ErrGenderNoFieldsToUpdate = errors.New("no fields to update")
	ErrGenderAlreadyExists    = errors.New("gender with the same name already exists")
)
