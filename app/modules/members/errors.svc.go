package members

import "errors"

var (
	ErrMemberNotFound              = errors.New("member not found")
	ErrMemberInvalidID             = errors.New("invalid member ID")
	ErrMemberInvalidGenderID       = errors.New("invalid member gender ID")
	ErrMemberInvalidPrefixID       = errors.New("invalid member prefix ID")
	ErrMemberNoFieldsToUpdate      = errors.New("no fields to update")
	ErrMemberUsernameRequired      = errors.New("username is required")
	ErrMemberPasswordRequired      = errors.New("password is required")
	ErrMemberUsernameExists        = errors.New("username already exists")
	ErrMemberInvalidCredentials    = errors.New("invalid credentials")
	ErrMemberUnauthorized          = errors.New("unauthorized")
	ErrMemberAccountNotFound       = errors.New("member account not found")
	ErrMemberPasswordMismatch      = errors.New("password mismatch")
	ErrMemberInvalidCurrency       = errors.New("invalid currency")
	ErrMemberInvalidLanguage       = errors.New("invalid language")
	ErrMemberNoSettingsToUpdate    = errors.New("no settings to update")
	ErrMemberStorageNotConfigured  = errors.New("member storage is not configured")
	ErrMemberGoogleOAuthDisabled   = errors.New("google oauth is not configured")
	ErrMemberGoogleInvalidState    = errors.New("google oauth state is invalid")
	ErrMemberGoogleAuthFailed      = errors.New("google oauth failed")
	ErrMemberNotificationInvalidID = errors.New("invalid member notification ID")
	ErrMemberNotificationNotFound  = errors.New("member notification not found")
)
