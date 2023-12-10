package bot

import "errors"

var (
	errUserInfoNotFound  = errors.New("error user info not found")
	errSendingSignUpLink = errors.New("error sending sing up link")
	errInvalidPetForm    = errors.New("error invalid pet form")
	errMissingPetField   = errors.New("missing field")
)
