package hooks

import (
	"context"

	"github.com/pkg/errors"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func ErrorLoggerHooks(log *zap.Logger) *twirp.ServerHooks {
	hooks := new(twirp.ServerHooks)

	hooks.Error = func(ctx context.Context, twerr twirp.Error) context.Context {

		if twerr.Code() == twirp.Internal {
			method, _ := twirp.MethodName(ctx)
			log.Error("internal error",
				zap.String("method", method),
				zap.Error(errors.Cause(twerr)),
			)
		}

		return ctx
	}

	return hooks
}
