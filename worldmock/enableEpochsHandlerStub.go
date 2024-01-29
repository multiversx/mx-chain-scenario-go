package worldmock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

var _ vmcommon.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	GetCurrentEpochCalled      func() uint32
	IsFlagDefinedCalled        func(flag core.EnableEpochFlag) bool
	IsFlagEnabledCalled        func(flag core.EnableEpochFlag) bool
	IsFlagEnabledInEpochCalled func(flag core.EnableEpochFlag, epoch uint32) bool
	GetActivationEpochCalled   func(flag core.EnableEpochFlag) uint32
}

// GetActivationEpoch -
func (stub *EnableEpochsHandlerStub) GetActivationEpoch(flag core.EnableEpochFlag) uint32 {
	if stub.GetActivationEpochCalled != nil {
		return stub.GetActivationEpochCalled(flag)
	}
	return 0
}

// IsFlagDefined -
func (stub *EnableEpochsHandlerStub) IsFlagDefined(flag core.EnableEpochFlag) bool {
	if stub.IsFlagDefinedCalled != nil {
		return stub.IsFlagDefinedCalled(flag)
	}
	return true
}

// IsFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFlagEnabled(flag core.EnableEpochFlag) bool {
	if stub.IsFlagEnabledCalled != nil {
		return stub.IsFlagEnabledCalled(flag)
	}

	return false
}

// IsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFlagEnabledInEpoch(flag core.EnableEpochFlag, epoch uint32) bool {
	if stub.IsFlagEnabledInEpochCalled != nil {
		return stub.IsFlagEnabledInEpochCalled(flag, epoch)
	}

	return false
}

// GetCurrentEpoch -
func (stub *EnableEpochsHandlerStub) GetCurrentEpoch() uint32 {
	if stub.GetCurrentEpochCalled != nil {
		return stub.GetCurrentEpochCalled()
	}

	return 0
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}

func EnableEpochsHandlerStubAllFlags() *EnableEpochsHandlerStub {
	return &EnableEpochsHandlerStub{
		GetCurrentEpochCalled: func() uint32 {
			return 0
		},
		IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
		IsFlagEnabledInEpochCalled: func(flag core.EnableEpochFlag, epoch uint32) bool {
			return true
		},
		GetActivationEpochCalled: func(flag core.EnableEpochFlag) uint32 {
			return 0
		},
	}
}
