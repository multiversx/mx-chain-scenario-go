package scenclibase

import (
	"errors"
	"fmt"
	"log"
	"os"

	mc "github.com/multiversx/mx-chain-scenario-go/scenario/io"

	cli "github.com/urfave/cli/v2"
)

// ScenariosCLI provides the functionality for any scenarios test executor.
func ScenariosCLI(version string, vmFlags CLIRunConfig) {
	app := cli.NewApp()
	app.Version = version
	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "print the tool version",
			Action: func(cCtx *cli.Context) error {
				fmt.Println(app.Version)
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "complete a task on the list",
			Flags: vmFlags.GetFlags(),
			Action: func(cCtx *cli.Context) error {
				args := cCtx.Args()
				if args.Len() != 1 {
					return errors.New("one path argument required to run scenarios")
				}
				path := cCtx.Args().First()

				return RunScenariosAtPath(path, vmFlags.ParseFlags(cCtx))
			},
		},
		{
			Name:  "fmt",
			Usage: "format all scenario files in a folder ( .scen.json / .step.json / .steps.json )",
			Action: func(cCtx *cli.Context) error {
				args := cCtx.Args()
				if args.Len() != 1 {
					return errors.New("one path argument required to format scenarios")
				}
				path := cCtx.Args().First()
				err := mc.FormatAllInFolder(path)
				return err
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
