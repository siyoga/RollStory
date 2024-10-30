package router

import (
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"

	"go.uber.org/zap"
)

func (r *router) defaultLogField(req domain.Request) []zap.Field {
	cmd := "default"
	if req.Command != nil {
		cmd = *req.Command
	}

	return []zap.Field{
		zap.Uint("id", uint(req.Id)), zap.String("command", cmd),
		zap.String("user_id", string(req.From.Id)),
	}
}

func (r *router) logRequest(req domain.Request, err *errors.Error) {
	var fields []zap.Field

	fields = append(r.defaultLogField(req), zap.Bool("success", true))

	if err != nil {
		fields = append(
			r.defaultLogField(req),
			zap.Bool("success", false),
			zap.String("reason", err.Reason),
			zap.NamedError("details", err.Details),
		)
	}

	r.logger.Info("request", fields...)
}
