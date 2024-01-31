package scenexec

import (
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
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

// VMBuilder defines VM initialization.
//
// The VM is not passed directly to the scenario executor because the scenario can override the gas schedule,
// which is required during initialization.
type VMBuilder interface {
	// NewMockWorld defines how the MockWorld is initialized.
	NewMockWorld() *worldmock.MockWorld

	// GasScheduleMapFromScenarios converts the gas schedule name from a scenario into an actual gas map.
	GasScheduleMapFromScenarios(scenGasSchedule scenmodel.GasSchedule) (worldmock.GasScheduleMap, error)

	// GetVMType returns the configured VM type, normally [5, 0].
	GetVMType() []byte

	// NewVM creates the execution VM host with references to the world mock and gas schedule.
	NewVM(world *worldmock.MockWorld, gasSchedule map[string]map[string]uint64) (VMInterface, error)
}
