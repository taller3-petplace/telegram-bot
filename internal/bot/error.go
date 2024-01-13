package bot

import "errors"

var (
	errUserInfoNotFound  = errors.New("error user info not found")
	errSendingSignUpLink = errors.New("error sending sing up link")
	errInvalidForm       = errors.New("error invalid form")
	errMissingFormField  = errors.New("missing form field")
)
