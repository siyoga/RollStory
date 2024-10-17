package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

type (
	Route struct {
		name string

		handler cmdHandler
	}

	Response struct {
		result domain.MessageResult
		code   int
		error  *errors.Error

		// pass the user ID of the person whose execution of the command has finished
		release *int64
	}

	cmdHandler func(ctx context.Context, userId int64, msg *domain.Message) Response
)

// Response
func NewSuccessResponse(res domain.MessageResult, code int, release *int64) Response {
	return Response{
		result:  res,
		code:    code,
		error:   nil,
		release: release,
	}
}

func NewErrResponse(e *errors.Error, release int64) Response {
	return Response{
		result:  domain.MessageResult{},
		code:    int(e.Code),
		error:   e,
		release: &release,
	}
}
