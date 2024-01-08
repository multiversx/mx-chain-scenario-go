package scencli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	mc "github.com/multiversx/mx-chain-scenario-go/controller"
	scenexec "github.com/multiversx/mx-chain-scenario-go/executor"
)

func run(path string, options CLIRunOptions) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	executor := scenexec.NewScenarioExecutor(options.VMBuilder)

	switch {
	case fi.IsDir():
		runner := mc.NewScenarioController(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunAllJSONScenariosInDirectory(
			path,
			"",
			".scen.json",
			[]string{},
			options.RunOptions)
	case strings.HasSuffix(path, ".scen.json"):
		runner := mc.NewScenarioController(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONScenario(path, options.RunOptions)
	default:
		err = errors.New("only directories and scenario files accepted as path")
	}

	// print result
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	return err
}
