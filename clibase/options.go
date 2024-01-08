package scenclibase

import (
	mc "github.com/multiversx/mx-chain-scenario-go/controller"
	scenexec "github.com/multiversx/mx-chain-scenario-go/executor"
	cli "github.com/urfave/cli/v2"
)

type CLIRunOptions struct {
	RunOptions *mc.RunScenarioOptions
	VMBuilder  scenexec.VMBuilder
}

type CLIRunConfig interface {
	GetFlags() []cli.Flag
	ParseFlags(cCtx *cli.Context) CLIRunOptions
}
