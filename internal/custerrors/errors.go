package custerrors

import "errors"

var (
	ErrNoContent                         = errors.New("no content")
	ErrOrderAlreadyUploadedBySameUser    = errors.New("order already uploaded by same user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrInsufficientFunds                 = errors.New("insufficient funds")
	ErrInvalidAuth                       = errors.New("invalid auth")

	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrInvalidOrderNumber = errors.New("invalid order number format")
)
