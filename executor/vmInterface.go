package scenexec

import (
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// VMInterface is the VM interface as defined in vm-common,
// plus a few extra methods needed for running the scenarios
type VMInterface interface {
	vmcommon.VMExecutionHandler

	Reset()
	SetGasTracing(enableGasTracing bool)
	GetGasTrace() map[string]map[string][]uint64
}

// ScenarioVMBuilder defines VM initialization.
//
// The VM is not passed directly to the scenario executor because the scenario can override the gas schedule,
// which is required during initialization.
type ScenarioVMBuilder interface {
	// GasScheduleMapFromScenarios converts the gas schedule name from a scenario into an actual gas map.
	GasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (worldmock.GasScheduleMap, error)

	// NewVM creates the execution VM host with references to the world mock and gas schedule.
	NewVM(world *worldmock.MockWorld, gasSchedule map[string]map[string]uint64) (VMInterface, error)
}
