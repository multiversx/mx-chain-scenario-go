package scenexec

import (
	mj "github.com/multiversx/mx-chain-scenario-go/model"
)

// ExecuteStep executes an individual step from a scenario.
func (ae *ScenarioExecutor) ExecuteStep(generalStep mj.Step) error {
	err := error(nil)

	switch step := generalStep.(type) {
	case *mj.ExternalStepsStep:
		err = ae.ExecuteExternalStep(step)
		length := len(ae.scenarioTraceGas)
		ae.scenarioTraceGas = ae.scenarioTraceGas[:length-1]
		return err
	case *mj.SetStateStep:
		err = ae.ExecuteSetStateStep(step)
	case *mj.CheckStateStep:
		err = ae.ExecuteCheckStateStep(step)
	case *mj.TxStep:
		_, err = ae.ExecuteTxStep(step)
	case *mj.DumpStateStep:
		err = ae.DumpWorld()
	}

	logGasTrace(ae)

	return err
}
