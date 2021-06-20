package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/google/subcommands"
	"github.com/suzuito/blog1-gcf"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/setting"
)

type runLocalCmd struct {
	gdeps *inject.GlobalDepends
	env   *setting.Environment
}

func newRunLocalCmd(gdeps *inject.GlobalDepends, env *setting.Environment) *runLocalCmd {
	return &runLocalCmd{gdeps: gdeps, env: env}
}

func (c *runLocalCmd) Name() string { return "run-local" }
func (c *runLocalCmd) Synopsis() string {
	return "GCFをローカルにて起動する\n"
}
func (c *runLocalCmd) Usage() string {
	return c.Synopsis()
}

func (c *runLocalCmd) SetFlags(f *flag.FlagSet) {
}

func (c *runLocalCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := blog1.RegisterLocalRunner(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return subcommands.ExitFailure
	}
	if err := funcframework.Start("8080"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
