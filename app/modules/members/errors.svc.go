package members

import "errors"

var (
	ErrMemberNotFound           = errors.New("member not found")
	ErrMemberInvalidID          = errors.New("invalid member ID")
	ErrMemberInvalidGenderID    = errors.New("invalid member gender ID")
	ErrMemberInvalidPrefixID    = errors.New("invalid member prefix ID")
	ErrMemberNoFieldsToUpdate   = errors.New("no fields to update")
	ErrMemberUsernameRequired   = errors.New("username is required")
	ErrMemberPasswordRequired   = errors.New("password is required")
	ErrMemberUsernameExists     = errors.New("username already exists")
	ErrMemberInvalidCredentials = errors.New("invalid credentials")
	ErrMemberUnauthorized       = errors.New("unauthorized")
	ErrMemberAccountNotFound    = errors.New("member account not found")
	ErrMemberPasswordMismatch   = errors.New("password mismatch")
)
