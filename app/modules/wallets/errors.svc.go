package wallets

import "errors"

var (
	ErrWalletNotFound         = errors.New("wallet not found")
	ErrWalletInvalidID        = errors.New("invalid wallet ID")
	ErrWalletInvalidMemberID  = errors.New("invalid wallet member ID")
	ErrWalletNameRequired     = errors.New("wallet name is required")
	ErrWalletBalanceInvalid   = errors.New("wallet balance must be non-negative")
	ErrWalletNoFieldsToUpdate = errors.New("no fields to update")
)
