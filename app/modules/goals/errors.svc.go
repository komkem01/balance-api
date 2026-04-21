package goals

import "errors"

var (
	ErrGoalInvalidID              = errors.New("invalid goal ID")
	ErrGoalNameRequired           = errors.New("goal name is required")
	ErrGoalNotFound               = errors.New("goal not found")
	ErrGoalInvalidMemberID        = errors.New("invalid goal member ID")
	ErrGoalStartDateInvalid       = errors.New("invalid goal start date")
	ErrGoalTargetDateInvalid      = errors.New("invalid goal target date")
	ErrGoalNoFieldsToUpdate       = errors.New("no fields to update")
	ErrGoalTypeInvalid            = errors.New("invalid goal type")
	ErrGoalStatusInvalid          = errors.New("invalid goal status")
	ErrGoalTargetAmountInvalid    = errors.New("invalid goal target amount")
	ErrGoalSourceTypeInvalid      = errors.New("invalid goal source type")
	ErrGoalSourceIDRequired       = errors.New("goal source id is required")
	ErrGoalSourceIDInvalid        = errors.New("invalid goal source id")
	ErrGoalSourceMemberForbidden  = errors.New("goal source does not belong to member")
	ErrGoalDepositWalletInvalid   = errors.New("invalid goal deposit wallet id")
	ErrGoalDepositWalletForbidden = errors.New("goal deposit wallet does not belong to member")
)
