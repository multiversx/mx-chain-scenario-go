package scenexec

import (
	"fmt"

	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

func logGasTrace(ae *ScenarioExecutor) {
	if ae.PeekTraceGas() {
		scGasTrace := ae.vm.GetGasTrace()
		totalGasUsedByAPIs := 0
		for scAddress, gasTrace := range scGasTrace {
			fmt.Println("Gas Trace for: ", "SC Address", scAddress)
			for functionName, value := range gasTrace {
				totalGasUsed := uint64(0)
				for _, usedGas := range value {
					totalGasUsed += usedGas
				}
				fmt.Println("GasTrace: functionName:", functionName, ",  totalGasUsed:", totalGasUsed, ", numberOfCalls:", len(value))
				totalGasUsedByAPIs += int(totalGasUsed)
			}
			fmt.Println("TotalGasUsedByAPIs: ", totalGasUsedByAPIs)
		}
	}
}

func setGasTraceInMetering(ae *ScenarioExecutor, enable bool) {
	if enable && ae.PeekTraceGas() {
		ae.vm.SetGasTracing(true)
	} else {
		ae.vm.SetGasTracing(false)
	}
}

func setExternalStepGasTracing(ae *ScenarioExecutor, step *scenmodel.ExternalStepsStep) {
	switch step.TraceGas.ToInt() {
	case scenmodel.Undefined.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, ae.PeekTraceGas())
	case scenmodel.TrueValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, true)
	case scenmodel.FalseValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, false)
	}
}

func resetGasTracesIfNewTest(ae *ScenarioExecutor, scenario *scenmodel.Scenario) {
	if ae.vm == nil || scenario.IsNewTest {
		ae.scenarioTraceGas = make([]bool, 0)
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, scenario.TraceGas)
		scenario.IsNewTest = false
	}
}
