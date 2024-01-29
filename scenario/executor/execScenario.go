package scenexec

import (
	fr "github.com/multiversx/mx-chain-scenario-go/scenario/expression/fileresolver"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// RunScenario executes an individual test.
func (ae *ScenarioExecutor) RunScenario(scenario *scenmodel.Scenario, fileResolver fr.FileResolver) error {
	ae.fileResolver = fileResolver
	ae.checkGas = scenario.CheckGas
	resetGasTracesIfNewTest(ae, scenario)

	err := ae.InitVM(scenario.GasSchedule)
	if err != nil {
		return err
	}

	txIndex := 0
	for _, generalStep := range scenario.Steps {
		setGasTraceInMetering(ae, true)
		err := ae.ExecuteStep(generalStep)
		if err != nil {
			return err
		}
		setGasTraceInMetering(ae, false)
		txIndex++
	}

	return nil
}
