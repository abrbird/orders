package models

import (
	"github.com/pkg/errors"
)

var (
	UnknownErrorMessage       = "unknown error"
	NotFoundErrorMessage      = "not found error"
	AlreadyExistsErrorMessage = "already exists error"

	RetryErrorMessage       = "please retry"
	MaxAttemptsErrorMessage = "reached max attempts"

	BrokerSendErrorMessage = "can not send message"
)

type retryError struct {
	Inner   error
	Message string
}

func NewRetryError(err error) error {
	return retryError{err, RetryErrorMessage}
}

func (e retryError) Error() string {
	return wrapError(e.Inner, e.Message).Error()
}

func (e retryError) Unwrap() error {
	return e.Inner
}

type maxAttemptsError struct {
	Inner   error
	Message string
}

func NewMaxAttemptsError(err error) error {
	return maxAttemptsError{err, MaxAttemptsErrorMessage}
}

func (e maxAttemptsError) Error() string {
	return wrapError(e.Inner, e.Message).Error()
}

func (e maxAttemptsError) Unwrap() error {
	return e.Inner
}

func wrapError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return errors.Wrap(err, message)
}

func UnknownError(err error) error {
	return wrapError(err, UnknownErrorMessage)
}
func NotFoundError(err error) error {
	return wrapError(err, NotFoundErrorMessage)
}
func AlreadyExistsError(err error) error {
	return wrapError(err, AlreadyExistsErrorMessage)
}
func BrokerSendError(err error) error {
	return wrapError(err, BrokerSendErrorMessage)
}

var RetryError = NewRetryError(nil)
var MaxAttemptsError = NewMaxAttemptsError(nil)
