package courier

import "errors"

// service
var (
	ErrNameEmpty    = errors.New("name empty")
	ErrPhoneEmpty   = errors.New("phone empty")
	ErrStatusEmpty  = errors.New("status empty")
	ErrInvalidPhone = errors.New("invalid phone")
	ErrInvalidId    = errors.New("invalid id")
	ErrEmptyRequest = errors.New("empty request")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidName = errors.New("invalid name")
)
// repo
var (
	ErrPhoneExists = errors.New("courier with that phone exists")
	ErrIdNotFound  = errors.New("id not found")
)