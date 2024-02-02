package scenclibase

import (
	"errors"
	"fmt"
	"os"
	"strings"

	scenexec "github.com/multiversx/mx-chain-scenario-go/scenario/executor"
	scenio "github.com/multiversx/mx-chain-scenario-go/scenario/io"
	scenjparse "github.com/multiversx/mx-chain-scenario-go/scenario/json/parse"
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
	controller := &scenio.ScenarioController{
		Executor: executor,
		Parser: scenjparse.NewParser(
			scenio.NewDefaultFileResolver(),
			options.VMBuilder.GetVMType()),
	}

	switch {
	case fi.IsDir():
		err = controller.RunAllJSONScenariosInDirectory(
			path,
			"",
			".scen.json",
			[]string{},
			options.RunOptions)
	case strings.HasSuffix(path, ".scen.json"):
		err = controller.RunSingleJSONScenario(path, options.RunOptions)
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
