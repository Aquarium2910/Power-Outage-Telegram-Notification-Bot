package entity

import "errors"

var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
)
