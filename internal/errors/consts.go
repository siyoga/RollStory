package errors

import "errors"

var (
	ErrTelegramGetUpdates  = errors.New("telegram get updates error")
	ErrTelegramSendMessage = errors.New("telegram send message error")

	ErrGptCreateMessage = errors.New("failed to create message in gpt")
	ErrGptCreateRun     = errors.New("failed to create run in gpt")
	ErrGptTimeout       = errors.New("failed to make request in time range")
	ErrGptRetrieveRun   = errors.New("failed to retrieve run from gpt")
	ErrGptCreateThread  = errors.New("failed to create thread in gpt")

	ErrRouterNoDefaultRoute = errors.New("no default route in router")

	ErrRedisGetRaw  = errors.New("redis get error")
	ErrRedisSaveRaw = errors.New("redis save error")

	ErrAdapterRequestFailed = errors.New("failed to make request with adapter")

	ErrTokenNotExist = errors.New("token does not exists")

	ErrAuthNumberAssignmentFailed = errors.New("number assignment failed")
	ErrAuthParseTokenRaw          = errors.New("parse token failed")

	ErrMessageSendFailed = errors.New("game send failed")
)

var (
	HashPassword = &Error{Code: 500, Reason: "hash password error"}

	ErrServiceUnavailable = &Error{Code: 500, Reason: "service unavailable"}
	ErrHttpRequest        = &Error{Reason: "http request error"}

	ErrPermissionDenied   = &Error{Code: 403, Reason: "permission denied"}
	ErrParse              = &Error{Code: 400, Reason: "parse failed"}
	ErrInternal           = &Error{Code: 500, Reason: "internal error"}
	ErrConflict           = &Error{Code: 409, Reason: "conflict"}
	ErrMissingCredentials = &Error{Code: 401, Reason: "missing credentials"}
	ErrValidation         = &Error{Code: 400, Reason: "validation failed"}
	ErrTimeout            = &Error{Code: 504, Reason: "timeout"}
)
