package transactions

import "errors"

var (
	ErrTransactionNotFound          = errors.New("transaction not found")
	ErrTransactionInvalidID         = errors.New("invalid transaction ID")
	ErrTransactionInvalidWalletID   = errors.New("invalid transaction wallet ID")
	ErrTransactionInvalidCategoryID = errors.New("invalid transaction category ID")
	ErrTransactionTypeInvalid       = errors.New("transaction type is invalid")
	ErrTransactionAmountInvalid     = errors.New("transaction amount must be non-negative")
	ErrTransactionInsufficientFunds = errors.New("insufficient wallet balance")
	ErrTransactionNoFieldsToUpdate  = errors.New("no fields to update")
	ErrTransactionDateInvalid       = errors.New("transaction date is invalid")
	ErrTransactionRangeInvalid      = errors.New("transaction range is invalid")
	ErrTransactionImageInvalid      = errors.New("transaction image is invalid")
	ErrTransactionImageTooLarge     = errors.New("transaction image is too large")
	ErrTransactionImageUploadFailed = errors.New("transaction image upload failed")
)
