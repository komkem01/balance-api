package goals

import "errors"

var (
	ErrGoalInvalidID             = errors.New("invalid goal ID")
	ErrGoalNameRequired          = errors.New("goal name is required")
	ErrGoalNotFound              = errors.New("goal not found")
	ErrGoalInvalidMemberID       = errors.New("invalid goal member ID")
	ErrGoalNoFieldsToUpdate      = errors.New("no fields to update")
	ErrGoalTypeInvalid           = errors.New("invalid goal type")
	ErrGoalStatusInvalid         = errors.New("invalid goal status")
	ErrGoalTargetAmountInvalid   = errors.New("invalid goal target amount")
	ErrGoalSourceTypeInvalid     = errors.New("invalid goal source type")
	ErrGoalSourceIDRequired      = errors.New("goal source id is required")
	ErrGoalSourceMemberForbidden = errors.New("goal source does not belong to member")
)
