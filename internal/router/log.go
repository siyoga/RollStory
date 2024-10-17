package router

import (
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"

	"go.uber.org/zap"
)

func (r *router) defaultLogField(u domain.Update) []zap.Field {
	return []zap.Field{
		zap.Uint("id", uint(u.ID)), zap.String("command", u.Message.GetCommand()),
		zap.String("user_id", string(u.Message.From.Id)),
	}
}

func (r *router) logRequest(u domain.Update, err *errors.Error) {
	var fields []zap.Field

	fields = append(r.defaultLogField(u), zap.Bool("success", true))

	if err != nil {
		fields = append(
			r.defaultLogField(u),
			zap.Bool("success", false),
			zap.String("reason", err.Reason),
			zap.NamedError("details", err.Details),
		)
	}

	r.logger.Info("request", fields...)
}
