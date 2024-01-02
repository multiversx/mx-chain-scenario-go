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

type ScenarioVMBuilder interface {
	GasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (worldmock.GasScheduleMap, error)

	NewVM(world *worldmock.MockWorld, gasSchedule map[string]map[string]uint64) (VMInterface, error)
}
