package wallets

import "errors"

var (
	ErrWalletNotFound         = errors.New("wallet not found")
	ErrWalletInvalidID        = errors.New("invalid wallet ID")
	ErrWalletInvalidMemberID  = errors.New("invalid wallet member ID")
	ErrWalletNameRequired     = errors.New("wallet name is required")
	ErrWalletNoFieldsToUpdate = errors.New("no fields to update")
)
