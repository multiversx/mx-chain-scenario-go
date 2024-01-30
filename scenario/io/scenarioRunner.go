package scenio

import (
	fr "github.com/multiversx/mx-chain-scenario-go/scenario/expression/fileresolver"
	scenjparse "github.com/multiversx/mx-chain-scenario-go/scenario/json/parse"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// ScenarioRunner describes a component that can run a VM scenario.
type ScenarioRunner interface {
	// Reset clears state/world.
	Reset()

	// RunScenario executes the scenario and checks if it passed. Failure is signaled by returning an error.
	// The FileResolver helps with resolving external steps.
	// TODO: group into a "execution context" param.
	RunScenario(*scenmodel.Scenario, fr.FileResolver) error
}

// ScenarioController is a component that can run json scenarios, using a provided executor.
type ScenarioController struct {
	Executor    ScenarioRunner
	RunsNewTest bool
	Parser      scenjparse.Parser
}

// NewScenarioController creates new ScenarioController instance.
func NewScenarioController(executor ScenarioRunner, fileResolver fr.FileResolver) *ScenarioController {
	return &ScenarioController{
		Executor: executor,
		Parser:   scenjparse.NewParser(fileResolver),
	}
}

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
// Reexported here to avoid having all external packages importing the parser.
// DefaultFileResolver is in parse for local tests only.
func NewDefaultFileResolver() *fr.DefaultFileResolver {
	return fr.NewDefaultFileResolver()
}
