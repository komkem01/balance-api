package transactions

import "errors"

var (
	ErrTransactionNotFound          = errors.New("transaction not found")
	ErrTransactionInvalidID         = errors.New("invalid transaction ID")
	ErrTransactionInvalidWalletID   = errors.New("invalid transaction wallet ID")
	ErrTransactionInvalidCategoryID = errors.New("invalid transaction category ID")
	ErrTransactionTypeInvalid       = errors.New("transaction type is invalid")
	ErrTransactionNoFieldsToUpdate  = errors.New("no fields to update")
	ErrTransactionDateInvalid       = errors.New("transaction date is invalid")
)
