package apperrors

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
)

type ErrorOption func(e *appError)

type PublicErr interface {
	Error() string
	InternalError() error
	Cause() error
}

type appError struct {
	message      string
	cause        error
	internalMsg  string
	internalData []interface{}
}

var (
	_ error     = (*appError)(nil)
	_ PublicErr = (*appError)(nil)
)

func Error(message string, opts ...ErrorOption) *appError { //nolint:revive // don't nind returning unexported here
	err := &appError{
		message: message,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func (e *appError) InternalError() error {
	var msg string

	if e.internalMsg != "" {
		msg = e.internalMsg
	} else {
		msg = e.message
	}

	if e.cause != nil {
		return errors.Wrap(e.cause, msg, e.internalData...)
	}
	return errors.WithDetails(errors.New(msg), e.internalData...)
}

func (e *appError) Error() string {
	return e.message
}

func (e *appError) Unwrap() error {
	return e.cause
}

func (e *appError) Cause() error { return e.cause }

func Show(logger logging.Logger, parent fyne.Window, err PublicErr, onClosed func()) {
	level.Error(logger).Err("application error", err.InternalError())
	d := dialog.NewError(err, parent)
	d.SetOnClosed(func() {
		if onClosed != nil {
			onClosed()
		}
	})
	d.Show()
}

func WithCause(c error) ErrorOption {
	return func(err *appError) {
		if ae, ok := c.(PublicErr); ok {
			err.message = fmt.Sprintf("%s: %s", err.message, ae.Error())
		}
		err.cause = c
	}
}

func WithInternalMessage(msg string) ErrorOption {
	return func(err *appError) {
		err.internalMsg = msg
	}
}

func WithInternalData(kvs ...interface{}) ErrorOption {
	return func(err *appError) {
		err.internalData = append(err.internalData, kvs...)
	}
}
