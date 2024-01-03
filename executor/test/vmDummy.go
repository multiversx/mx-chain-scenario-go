package executortest

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	scenarioexec "github.com/multiversx/mx-chain-scenario-go/executor"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

var _ scenarioexec.VMInterface = (*DummyVM)(nil)
var _ scenarioexec.ScenarioVMBuilder = (*DummyVMBuilder)(nil)

// DummyVM is a VM stand-in that can never be called.
// Used for tests that do not require a VM.
type DummyVM struct{}

// RunSmartContractCreate -
func (*DummyVM) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	return nil, errors.New("cannot call the DummyVM")
}

// RunSmartContractCall -
func (*DummyVM) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	return nil, errors.New("cannot call the DummyVM")
}

// GasScheduleChange -
func (*DummyVM) GasScheduleChange(newGasSchedule map[string]map[string]uint64) {}

// GetVersion -
func (*DummyVM) GetVersion() string {
	return ""
}

// IsInterfaceNil -
func (*DummyVM) IsInterfaceNil() bool {
	return false
}

// Close -
func (*DummyVM) Close() error {
	return nil
}

// Reset -
func (*DummyVM) Reset() {}

// SetGasTracing -
func (*DummyVM) SetGasTracing(enableGasTracing bool) {}

// GetGasTrace -
func (*DummyVM) GetGasTrace() map[string]map[string][]uint64 {
	return make(map[string]map[string][]uint64)
}

// DummyVMBuilder is the builder for a DummyVM.
// Also provides a minimal gas schedule for running the builtin functions.
// Used for tests that do not require a VM.
type DummyVMBuilder struct{}

func (*DummyVMBuilder) GasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (worldmock.GasScheduleMap, error) {
	gasMap := make(map[string]map[string]uint64)
	fillGasMapInternal(gasMap, 1)
	return gasMap, nil
}

func (*DummyVMBuilder) NewVM(world *worldmock.MockWorld, gasSchedule map[string]map[string]uint64) (scenarioexec.VMInterface, error) {
	return &DummyVM{}, nil
}

func fillGasMapInternal(gasMap map[string]map[string]uint64, value uint64) map[string]map[string]uint64 {
	gasMap[core.BaseOperationCostString] = fillGasMapBaseOperationCosts(value)
	gasMap[core.BuiltInCostString] = fillGasMapBuiltInCosts(value)

	return gasMap
}

func fillGasMapBaseOperationCosts(value uint64) map[string]uint64 {
	gasMap := make(map[string]uint64)
	gasMap["StorePerByte"] = value
	gasMap["DataCopyPerByte"] = value
	gasMap["ReleasePerByte"] = value
	gasMap["PersistPerByte"] = value
	gasMap["CompilePerByte"] = value
	gasMap["AoTPreparePerByte"] = value
	gasMap["GetCode"] = value
	return gasMap
}

func fillGasMapBuiltInCosts(value uint64) map[string]uint64 {
	gasMap := make(map[string]uint64)
	gasMap["ChangeOwnerAddress"] = value
	gasMap["ClaimDeveloperRewards"] = value
	gasMap["SaveUserName"] = value
	gasMap["SaveKeyValue"] = value
	gasMap["ESDTTransfer"] = value
	gasMap["ESDTBurn"] = value
	gasMap["ESDTLocalMint"] = value
	gasMap["ESDTLocalBurn"] = value
	gasMap["ESDTNFTCreate"] = value
	gasMap["ESDTNFTAddQuantity"] = value
	gasMap["ESDTNFTBurn"] = value
	gasMap["ESDTNFTTransfer"] = value
	gasMap["ESDTNFTChangeCreateOwner"] = value
	gasMap["ESDTNFTAddUri"] = value
	gasMap["ESDTNFTUpdateAttributes"] = value
	gasMap["ESDTNFTMultiTransfer"] = value
	gasMap["SetGuardian"] = value
	gasMap["GuardAccount"] = value
	gasMap["UnGuardAccount"] = value
	gasMap["TrieLoadPerNode"] = value
	gasMap["TrieStorePerNode"] = value

	return gasMap
}
