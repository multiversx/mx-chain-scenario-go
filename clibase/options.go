package scenclibase

import (
	scenexec "github.com/multiversx/mx-chain-scenario-go/scenario/executor"
	scenio "github.com/multiversx/mx-chain-scenario-go/scenario/io"

	cli "github.com/urfave/cli/v2"
)

// CLIRunOptions are all the options needed to run scenarios in a directory.
type CLIRunOptions struct {
	RunOptions *scenio.RunScenarioOptions
	VMBuilder  scenexec.VMBuilder
}

// CLIRunConfig prepares and interprets CLI flags required to run scenarios at a path.
type CLIRunConfig interface {
	GetFlags() []cli.Flag
	ParseFlags(cCtx *cli.Context) CLIRunOptions
}
