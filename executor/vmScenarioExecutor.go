package scenexec

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	mc "github.com/multiversx/mx-chain-scenario-go/controller"
	er "github.com/multiversx/mx-chain-scenario-go/expression/reconstructor"
	fr "github.com/multiversx/mx-chain-scenario-go/fileresolver"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

var log = logger.GetOrCreate("scenario-exec")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ScenarioExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type ScenarioExecutor struct {
	World             *worldmock.MockWorld
	vmBuilder         VMBuilder
	vm                VMInterface
	checkGas          bool
	scenarioTraceGas  []bool
	fileResolver      fr.FileResolver
	exprReconstructor er.ExprReconstructor
}

var _ mc.ScenarioRunner = (*ScenarioExecutor)(nil)

// NewScenarioExecutor prepares a new VMTestExecutor instance.
func NewScenarioExecutor(vmBuilder VMBuilder) *ScenarioExecutor {
	world := worldmock.NewMockWorld()

	return &ScenarioExecutor{
		World:             world,
		vm:                nil,
		vmBuilder:         vmBuilder,
		checkGas:          true,
		scenarioTraceGas:  make([]bool, 0),
		fileResolver:      nil,
		exprReconstructor: er.ExprReconstructor{},
	}
}

// InitVM will initialize the VM and the builtin function container.
// Does nothing if the VM is already initialized.
func (ae *ScenarioExecutor) InitVM(scenGasSchedule mj.GasSchedule) error {
	if ae.vm != nil {
		return nil
	}

	gasSchedule, err := ae.vmBuilder.GasScheduleMapFromScenarios(scenGasSchedule)
	if err != nil {
		return err
	}

	err = ae.World.InitBuiltinFunctions(gasSchedule)
	if err != nil {
		return err
	}

	ae.vm, err = ae.vmBuilder.NewVM(ae.World, gasSchedule)

	return err
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ScenarioExecutor) GetVM() vmcommon.VMExecutionHandler {
	return ae.vm
}

// Reset clears state/world.
// Is called in RunAllJSONScenariosInDirectory, but not in RunSingleJSONScenario.
func (ae *ScenarioExecutor) Reset() {
	if !check.IfNil(ae.vm) {
		ae.vm.Reset()
	}
	ae.World.Clear()
}

// Close will simply close the VM
func (ae *ScenarioExecutor) Close() {
	if !check.IfNil(ae.vm) {
		ae.vm.Reset()
	}
}

// PeekTraceGas returns the last position from the scenarioTraceGas, if existing
func (ae *ScenarioExecutor) PeekTraceGas() bool {
	length := len(ae.scenarioTraceGas)
	if length != 0 {
		return ae.scenarioTraceGas[length-1]
	}
	return false
}
