package v1

import "errors"

var (
	errUserAlreadyExists = errors.New("user with such email already exists")
	errUserNotFound      = errors.New("user doesn't exists")

	noId            = "no id"
	noCode          = "code is empty"
	errInvalidInput = "invalid input body"
)
