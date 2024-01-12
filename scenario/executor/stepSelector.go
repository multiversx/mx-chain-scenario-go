package scenexec

import (
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// ExecuteStep executes an individual step from a scenario.
func (ae *ScenarioExecutor) ExecuteStep(generalStep scenmodel.Step) error {
	err := error(nil)

	switch step := generalStep.(type) {
	case *scenmodel.ExternalStepsStep:
		err = ae.ExecuteExternalStep(step)
		length := len(ae.scenarioTraceGas)
		ae.scenarioTraceGas = ae.scenarioTraceGas[:length-1]
		return err
	case *scenmodel.SetStateStep:
		err = ae.ExecuteSetStateStep(step)
	case *scenmodel.CheckStateStep:
		err = ae.ExecuteCheckStateStep(step)
	case *scenmodel.TxStep:
		_, err = ae.ExecuteTxStep(step)
	case *scenmodel.DumpStateStep:
		err = ae.DumpWorld()
	}

	logGasTrace(ae)

	return err
}
