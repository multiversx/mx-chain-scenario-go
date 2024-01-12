package scenclibase

import (
	"errors"
	"fmt"
	"os"
	"strings"

	scenexec "github.com/multiversx/mx-chain-scenario-go/scenario/executor"
	mc "github.com/multiversx/mx-chain-scenario-go/scenario/io"
)

// RunScenariosAtPath runs either;
// - all scenarios in folder if path is a directory
// - single scenario given as path.
func RunScenariosAtPath(path string, options CLIRunOptions) error {
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
