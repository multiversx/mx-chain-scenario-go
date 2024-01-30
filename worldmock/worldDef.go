package worldmock

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// DefaultHasher is an exposed value to use in tests
var DefaultHasher = blake2b.NewBlake2b()

// DefaultVMType is an exposed value to use in tests
var DefaultVMType = []byte{0xF, 0xF}

// NewAddressMock allows tests to specify what new addresses to generate
type NewAddressMock struct {
	CreatorAddress []byte
	CreatorNonce   uint64
	NewAddress     []byte
}

// BlockInfo contains metadata about a mocked block
type BlockInfo struct {
	BlockTimestamp uint64
	BlockNonce     uint64
	BlockRound     uint64
	BlockEpoch     uint32
	RandomSeed     *[48]byte
}

// GetRandomSeedSlice retrieves the configured random seed or a slice of zeros.
// Always 48 bytes long, never nil.
func (bi *BlockInfo) GetRandomSeedSlice() []byte {
	if bi.RandomSeed == nil {
		bi.RandomSeed = &[48]byte{}
	}
	return bi.RandomSeed[:]
}

// MockWorld provides a mock representation of the blockchain to be used in VM tests.
type MockWorld struct {
	SelfShardID                uint32
	VMType                     []byte
	AcctMap                    AccountMap
	AccountsAdapter            vmcommon.AccountsAdapter
	PreviousBlockInfo          *BlockInfo
	CurrentBlockInfo           *BlockInfo
	Blockhashes                [][]byte
	NewAddressMocks            []*NewAddressMock
	StateRootHash              []byte
	Err                        error
	LastCreatedContractAddress []byte
	CompiledCode               map[string][]byte
	BuiltinFuncs               *BuiltinFunctionsWrapper
	IsPausedValue              bool
	IsLimitedTransferValue     bool
	GuardedAccountHandler      vmcommon.GuardedAccountHandler
	ProvidedBlockchainHook     vmcommon.BlockchainHook
	EnableEpochsHandler        vmcommon.EnableEpochsHandler
	OtherVMOutputMap           map[string]*vmcommon.VMOutput
}

// NewMockWorld creates a new MockWorld instance
func NewMockWorld() *MockWorld {
	accountMap := NewAccountMap()
	world := &MockWorld{
		SelfShardID:         0,
		AcctMap:             accountMap,
		AccountsAdapter:     nil,
		PreviousBlockInfo:   nil,
		CurrentBlockInfo:    nil,
		Blockhashes:         nil,
		NewAddressMocks:     nil,
		CompiledCode:        make(map[string][]byte),
		BuiltinFuncs:        nil,
		EnableEpochsHandler: EnableEpochsHandlerStubAllFlags(),
		OtherVMOutputMap:    make(map[string]*vmcommon.VMOutput),
	}
	world.AccountsAdapter = NewMockAccountsAdapter(world)
	world.GuardedAccountHandler = NewMockGuardedAccountHandler()

	return world
}

// SetProvidedBlockchainHook -
func (b *MockWorld) SetProvidedBlockchainHook(bh vmcommon.BlockchainHook) {
	b.ProvidedBlockchainHook = bh
}

// InitBuiltinFunctions initializes the inner BuiltinFunctionsWrapper, required
// for calling builtin functions.
func (b *MockWorld) InitBuiltinFunctions(gasMap GasScheduleMap) error {
	wrapper, err := NewBuiltinFunctionsWrapper(b, gasMap)
	if err != nil {
		return err
	}

	b.BuiltinFuncs = wrapper
	return nil
}

// Clear resets all mock data between tests.
func (b *MockWorld) Clear() {
	b.AcctMap = NewAccountMap()
	b.AccountsAdapter = NewMockAccountsAdapter(b)
	b.PreviousBlockInfo = nil
	b.CurrentBlockInfo = nil
	b.Blockhashes = nil
	b.NewAddressMocks = nil
	b.CompiledCode = make(map[string][]byte)
}

// SetCurrentBlockHash -
func (b *MockWorld) SetCurrentBlockHash(blockHash []byte) {
	if b.CurrentBlockInfo == nil {
		b.CurrentBlockInfo = &BlockInfo{}
	}
	b.Blockhashes = [][]byte{blockHash}
}

// NumberOfShards -
func (b *MockWorld) NumberOfShards() uint32 {
	maxShardID := uint32(0)
	for _, account := range b.AcctMap {
		if account.ShardID > maxShardID {
			maxShardID = account.ShardID
		}
	}

	return maxShardID + 1
}

// ComputeId -
func (b *MockWorld) ComputeId(address []byte) uint32 {
	return b.AcctMap.GetAccount(address).ShardID
}

// SelfId -
func (b *MockWorld) SelfId() uint32 {
	return b.SelfShardID
}

// SameShard -
func (b *MockWorld) SameShard(firstAddress []byte, secondAddress []byte) bool {
	firstAccount := b.AcctMap.GetAccount(firstAddress)
	secondAccount := b.AcctMap.GetAccount(secondAddress)
	return firstAccount.ShardID == secondAccount.ShardID
}

// CommunicationIdentifier -
func (b *MockWorld) CommunicationIdentifier(destShardID uint32) string {
	return fmt.Sprintf("commID-dest-%d", destShardID)
}

// GetSnapshot -
func (b *MockWorld) GetSnapshot() int {
	b.CreateStateBackup()
	return b.AccountsAdapter.JournalLen()
}

// RevertToSnapshot -
func (b *MockWorld) RevertToSnapshot(snapshot int) error {
	return b.AccountsAdapter.RevertToSnapshot(snapshot)
}

// ExecuteSmartContractCallOnOtherVM -
func (b *MockWorld) ExecuteSmartContractCallOnOtherVM(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	vmType, err := vmcommon.ParseVMTypeFromContractAddress(input.RecipientAddr)
	if err != nil {
		return nil, err
	}
	vmOutput := b.OtherVMOutputMap[string(vmType)]
	if vmOutput == nil {
		return &vmcommon.VMOutput{}, nil
	}

	return vmOutput, nil
}
