package logger

import (
	"fmt"

	"github.com/siyoga/rollstory/internal/errors"

	"go.uber.org/zap"
)

type (
	Logger interface {
		Named(name string) Logger

		Zap() *zap.Logger

		Error(err, reason error)
		DbError(err, reason error, query string) error
		AdapterError(err, reason error, request string) error

		ServiceError(err *errors.Error) *errors.Error
		// ServiceTxError(err error) *errors.Error
		ServiceDatabaseError(err error) *errors.Error

		Info(text string, fields ...zap.Field)
		Panic(text string, reason error)
	}

	logger struct {
		moduleName string
		log        *zap.Logger
	}
)

func NewLogger(l *zap.Logger, name string, moduleName string) Logger {
	return &logger{
		moduleName: moduleName,
		log:        l.Named(name),
	}
}

func (l *logger) Named(name string) Logger {
	return &logger{
		log: l.log.Named(name),
	}
}

func (l *logger) Zap() *zap.Logger {
	return l.log
}

func (l *logger) Panic(text string, reason error) {
	l.log.Panic(text, zap.NamedError("reason", reason))
}
func (l *logger) Info(text string, fields ...zap.Field) {
	l.log.Info(text, fields...)
}
func (l *logger) Warn(text string, fields ...zap.Field) {
	l.log.Warn(text, fields...)
}
func (l *logger) Error(err, reason error) {
	if reason == nil {
		reason = err
	}

	l.log.Error(l.genError(err, reason.Error()))
}
func (l *logger) DbError(err, reason error, query string) error {
	if reason == nil {
		reason = err
	}

	l.log.Error(l.genError(fmt.Errorf("%s \n %s", err.Error(), query), reason.Error()))
	return err
}

func (l *logger) AdapterError(err, reason error, request string) error {
	if reason == nil {
		reason = err
	}

	l.log.Error(l.genError(fmt.Errorf("%s \n %s", err.Error(), request), reason.Error()))
	return err
}

func (l *logger) ServiceDatabaseError(err error) *errors.Error {
	return l.ServiceError(errors.DatabaseError(err))
}

func (l *logger) ServiceError(err *errors.Error) *errors.Error {
	e := fmt.Errorf(err.Reason)
	if err.Details != nil {
		e = err.Details
	}

	l.log.Error(l.genError(e, err.Reason))
	return err
}

// func (l *logger) ServiceTxError(err error) *errors.Error {
// 	l.log.Error(l.genError(err, errors.ErrPostgresTx.Error()))
// 	return errors.DatabaseError(err)
// }
