package blog1

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/setting"
)

var closeFunc func()
var env *setting.Environment
var gdeps *inject.GlobalDepends

func init() {
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "zerolog_timestamp"
	ctxGlobal := context.Background()
	var err error
	env, err = setting.NewEnvironment()
	if err != nil {
		log.Error().AnErr("message", err).Send()
		return
	}
	gdeps, closeFunc, err = inject.NewGlobalDepends(ctxGlobal, env)
	if err != nil {
		log.Error().AnErr("message", err).Send()
		return
	}
}
