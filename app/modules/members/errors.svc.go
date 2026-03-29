package members

import "errors"

var (
	ErrMemberNotFound         = errors.New("member not found")
	ErrMemberInvalidID        = errors.New("invalid member ID")
	ErrMemberInvalidGenderID  = errors.New("invalid member gender ID")
	ErrMemberInvalidPrefixID  = errors.New("invalid member prefix ID")
	ErrMemberNoFieldsToUpdate = errors.New("no fields to update")
)
