package memberaccounts

import "errors"

var (
	ErrMemberAccountNotFound         = errors.New("member account not found")
	ErrMemberAccountInvalidID        = errors.New("invalid member account ID")
	ErrMemberAccountInvalidMemberID  = errors.New("invalid member account member ID")
	ErrMemberAccountNoFieldsToUpdate = errors.New("no fields to update")
)
