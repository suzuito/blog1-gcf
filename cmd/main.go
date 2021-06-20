package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/rs/zerolog/log"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/setting"
)

func main() {
	ctx := context.Background()
	env, err := setting.NewEnvironment()
	if err != nil {
		log.Error().AnErr("message", err).Send()
		return
	}
	gdeps, closeFunc, err := inject.NewGlobalDepends(ctx, env)
	if err != nil {
		log.Error().AnErr("message", err).Send()
		return
	}
	defer closeFunc()
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(newRunLocalCmd(gdeps, env), "")
	flag.Parse()
	os.Exit(int(subcommands.Execute(ctx)))
}
