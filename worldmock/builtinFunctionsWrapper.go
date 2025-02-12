package worldmock

import (
	"bytes"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
)

// WorldMarshalizer is the global marshalizer to be used by the components of
// the BuiltinFunctionsWrapper.
var WorldMarshalizer = &marshal.GogoProtoMarshalizer{}

// GasScheduleMap (alias) is the map for gas schedule
type GasScheduleMap = map[string]map[string]uint64

// BuiltinFunctionsWrapper manages and initializes a BuiltInFunctionContainer
// along with its dependencies
type BuiltinFunctionsWrapper struct {
	Container       vmcommon.BuiltInFunctionContainer
	MapDNSAddresses map[string]struct{}
	World           *MockWorld
	Marshalizer     vmcommon.Marshalizer
}

// NewBuiltinFunctionsWrapper creates a new BuiltinFunctionsWrapper with
// default dependencies.
func NewBuiltinFunctionsWrapper(
	world *MockWorld,
	gasMap GasScheduleMap,
) (*BuiltinFunctionsWrapper, error) {

	dnsMap := makeDNSAddresses(numDNSAddresses)

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasMap:                           gasMap,
		MapDNSAddresses:                  dnsMap,
		MapDNSV2Addresses:                dnsMap,
		Marshalizer:                      WorldMarshalizer,
		Accounts:                         world.AccountsAdapter,
		ShardCoordinator:                 world,
		EnableEpochsHandler:              world.EnableEpochsHandler,
		GuardedAccountHandler:            world.GuardedAccountHandler,
		MaxNumOfAddressesForTransferRole: 100,
	}

	builtinFuncFactory, err := builtInFunctions.NewBuiltInFunctionsCreator(argsBuiltIn)
	if err != nil {
		return nil, err
	}

	err = builtinFuncFactory.CreateBuiltInFunctionContainer()
	if err != nil {
		return nil, err
	}

	err = builtinFuncFactory.SetPayableHandler(world)
	if err != nil {
		return nil, err
	}

	builtinFuncsWrapper := &BuiltinFunctionsWrapper{
		Container:       builtinFuncFactory.BuiltInFunctionContainer(),
		MapDNSAddresses: argsBuiltIn.MapDNSAddresses,
		World:           world,
	}

	return builtinFuncsWrapper, nil
}

func (bf *BuiltinFunctionsWrapper) ensureAccountExists(address []byte) vmcommon.UserAccountHandler {
	account := bf.World.AcctMap.GetAccount(address)
	if account == nil {
		account = bf.World.AcctMap.CreateAccount(address, bf.World)
	}
	return account
}

// ProcessBuiltInFunction delegates the execution of a real builtin function to
// the inner BuiltInFunctionContainer.
func (bf *BuiltinFunctionsWrapper) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	caller := bf.getAccountSharded(input.CallerAddr)

	_ = bf.ensureAccountExists(input.RecipientAddr)
	recipient := bf.getAccountSharded(input.RecipientAddr)

	function, err := bf.Container.Get(input.Function)
	if err != nil {
		return nil, err
	}

	vmOutput, err := function.ProcessBuiltinFunction(caller, recipient, input)
	if err != nil {
		return nil, err
	}

	if !check.IfNil(caller) {
		err = bf.World.AccountsAdapter.SaveAccount(caller)
		if err != nil {
			return nil, err
		}
	}

	if !check.IfNil(recipient) && !bytes.Equal(input.CallerAddr, input.RecipientAddr) {
		err = bf.World.AccountsAdapter.SaveAccount(recipient)
		if err != nil {
			return nil, err
		}
	}

	return vmOutput, nil
}

// GetBuiltinFunctionNames returns the list of defined builtin-in functions.
func (bf *BuiltinFunctionsWrapper) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return bf.Container.Keys()
}

// TODO change AccountMap to support this instead
func (bf *BuiltinFunctionsWrapper) getAccountSharded(address []byte) vmcommon.UserAccountHandler {
	accountShard := bf.World.ComputeId(address)
	if accountShard != bf.World.SelfId() {
		return nil
	}
	return bf.World.AcctMap.GetAccount(address)
}
