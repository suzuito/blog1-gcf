package blog1

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/setting"
)

var closeFunc func()
var env *setting.Environment
var gdeps *inject.GlobalDepends

func init() {
	ctxGlobal := context.Background()
	var err error
	env, err = setting.NewEnvironment()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	gdeps, closeFunc, err = inject.NewGlobalDepends(ctxGlobal, env)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
}
