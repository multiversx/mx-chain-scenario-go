package worldmock

import (
	"bytes"
	"errors"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-crypto-go/address"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// ErrOperationNotPermitted indicates an operation rejected due to insufficient
// permissions.
var ErrOperationNotPermitted = errors.New("operation not permitted")

// ErrInvalidAddressLength indicates an incorrect length given for an address.
var ErrInvalidAddressLength = errors.New("invalid address length")

// Account holds the account info
type Account struct {
	Exists          bool
	Address         []byte
	Nonce           uint64
	Balance         *big.Int
	BalanceDelta    *big.Int
	Storage         map[string][]byte
	RootHash        []byte
	Code            []byte
	CodeHash        []byte
	CodeMetadata    []byte
	OwnerAddress    []byte
	AsyncCallData   string
	Username        []byte
	DeveloperReward *big.Int
	ShardID         uint32
	IsSmartContract bool
	MainIdentifier  core.AddressIdentifier
	Aliases         map[core.AddressIdentifier][]byte
	MockWorld       *MockWorld
}

var storageDefaultValue = make([]byte, 0)

// StorageValue yields the storage value for key, default 0
func (a *Account) StorageValue(key string) []byte {
	value, found := a.Storage[key]
	if !found {
		return storageDefaultValue
	}
	return value
}

// SetCodeAndMetadata changes the account code, as well as all fields depending on it:
// CodeHash, IsSmartContract, CodeMetadata.
// The code metadata must be given explicitly.
func (a *Account) SetCodeAndMetadata(code []byte, codeMetadata *vmcommon.CodeMetadata) {
	a.Code = code
	hash := DefaultHasher.Compute(string(a.Code))

	a.CodeHash = hash
	a.IsSmartContract = true
	a.CodeMetadata = codeMetadata.ToBytes()
}

// AddressBytes -
func (a *Account) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// GetCode -
func (a *Account) GetCode() []byte {
	return a.Code
}

// GetCodeMetadata -
func (a *Account) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *Account) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash -
func (a *Account) GetRootHash() []byte {
	return a.RootHash
}

// GetBalance -
func (a *Account) GetBalance() *big.Int {
	return a.Balance
}

// SetBalance -
func (a *Account) SetBalance(balance int64) {
	a.Balance = big.NewInt(balance)
}

// GetDeveloperReward -
func (a *Account) GetDeveloperReward() *big.Int {
	return a.DeveloperReward
}

// GetOwnerAddress -
func (a *Account) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetUserName -
func (a *Account) GetUserName() []byte {
	return a.Username
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}

// SetCode -
func (a *Account) SetCode(code []byte) {
	a.Code = code
	a.CodeHash = DefaultHasher.Compute(string(code))
	a.IsSmartContract = true
}

// SetCodeMetadata -
func (a *Account) SetCodeMetadata(codeMetadata []byte) {
	a.CodeMetadata = codeMetadata
}

// SetCodeHash -
func (a *Account) SetCodeHash(hash []byte) {
	a.CodeHash = hash
}

// SetRootHash -
func (a *Account) SetRootHash(hash []byte) {
	a.RootHash = hash
}

// AccountDataHandler -
func (a *Account) AccountDataHandler() vmcommon.AccountDataHandler {
	return a
}

// AddToBalance -
func (a *Account) AddToBalance(value *big.Int) error {
	newBalance := big.NewInt(0).Add(a.Balance, value)
	if newBalance.Cmp(zero) < 0 {
		return ErrInsufficientFunds
	}

	a.Balance = newBalance
	return nil
}

// SubFromBalance -
func (a *Account) SubFromBalance(value *big.Int) error {
	newBalance := big.NewInt(0).Sub(a.Balance, value)
	if newBalance.Cmp(zero) < 0 {
		return ErrInsufficientFunds
	}

	a.Balance = newBalance
	return nil
}

// ClaimDeveloperRewards -
func (a *Account) ClaimDeveloperRewards(sender []byte) (*big.Int, error) {
	if !bytes.Equal(sender, a.OwnerAddress) {
		return nil, ErrOperationNotPermitted
	}

	oldValue := big.NewInt(0).Set(a.DeveloperReward)
	a.DeveloperReward = big.NewInt(0)

	return oldValue, nil
}

// AddToDeveloperReward -
func (a *Account) AddToDeveloperReward(value *big.Int) {
	a.DeveloperReward = big.NewInt(0).Add(a.DeveloperReward, value)
}

// ChangeOwnerAddress -
func (a *Account) ChangeOwnerAddress(sender []byte, newAddress []byte) error {
	if !bytes.Equal(sender, a.OwnerAddress) {
		return ErrOperationNotPermitted
	}
	if len(newAddress) != len(a.Address) {
		return ErrInvalidAddressLength
	}

	a.OwnerAddress = newAddress

	return nil
}

// SetOwnerAddress -
func (a *Account) SetOwnerAddress(address []byte) {
	a.OwnerAddress = address
}

// SetUserName -
func (a *Account) SetUserName(userName []byte) {
	a.Username = make([]byte, len(userName))
	copy(a.Username, userName)
}

// IncreaseNonce -
func (a *Account) IncreaseNonce(nonce uint64) {
	a.Nonce += nonce
}

// RetrieveValue -
func (a *Account) RetrieveValue(key []byte) ([]byte, uint32, error) {
	return a.Storage[string(key)], 0, nil
}

// MigrateDataTrieLeaves -
func (a *Account) MigrateDataTrieLeaves(_ vmcommon.ArgsMigrateDataTrieLeaves) error {
	return nil
}

// SaveKeyValue -
func (a *Account) SaveKeyValue(key []byte, value []byte) error {
	a.Storage[string(key)] = value
	if a.MockWorld == nil {
		return ErrNilWorldMock
	}
	a.MockWorld.CreateStateBackup()
	return nil
}

// ClearDataCaches -
func (a *Account) ClearDataCaches() {
}

// DirtyData -
func (a *Account) DirtyData() map[string][]byte {
	return a.Storage
}

// Clone -
func (a *Account) Clone() *Account {
	return &Account{
		Exists:          a.Exists,
		Address:         a.Address,
		Nonce:           a.Nonce,
		Balance:         big.NewInt(0).Set(a.Balance),
		BalanceDelta:    big.NewInt(0).Set(a.BalanceDelta),
		Storage:         a.cloneStorage(),
		RootHash:        cloneBytes(a.RootHash),
		Code:            cloneBytes(a.Code),
		CodeHash:        cloneBytes(a.CodeHash),
		CodeMetadata:    cloneBytes(a.CodeMetadata),
		AsyncCallData:   a.AsyncCallData,
		OwnerAddress:    cloneBytes(a.OwnerAddress),
		Username:        cloneBytes(a.Username),
		DeveloperReward: big.NewInt(0).Set(a.DeveloperReward),
		ShardID:         a.ShardID,
		IsSmartContract: a.IsSmartContract,
		Aliases:         a.Aliases,
		MockWorld:       a.MockWorld,
	}
}

func (a *Account) cloneStorage() map[string][]byte {
	clone := make(map[string][]byte, len(a.Storage))
	for key, value := range a.Storage {
		clone[key] = cloneBytes(value)
	}

	return clone
}

func cloneBytes(b []byte) []byte {
	clone := make([]byte, len(b))
	copy(clone, b)
	return clone
}

// GetAlias -
func (a *Account) GetAlias(identifier core.AddressIdentifier) []byte {
	return a.Aliases[identifier]
}

// SetAlias -
func (a *Account) SetAlias(aliasAddress []byte, aliasIdentifier core.AddressIdentifier) {
	a.Aliases[aliasIdentifier] = aliasAddress
	a.MockWorld.AliasesMap[aliasIdentifier.BuildAddressIdentifier(aliasAddress)] = a.Address
}

func (a *Account) RequestMainAddressIdentifier(request *vmcommon.AddressRequest) core.AddressIdentifier {
	switch a.MainIdentifier {
	case core.InvalidAddressIdentifier:
		if request.SaveOnGenerate {
			a.MainIdentifier = request.SourceIdentifier
		}
		return request.SourceIdentifier
	default:
		return a.MainIdentifier
	}
}

func (a *Account) RequestMainAddress(request *vmcommon.AddressRequest) ([]byte, core.AddressIdentifier) {
	mainAddressIdentifier := a.RequestMainAddressIdentifier(request)

	switch mainAddressIdentifier {
	case request.SourceIdentifier:
		return request.SourceAddress, mainAddressIdentifier
	case core.MVXAddressIdentifier:
		return a.AddressBytes(), mainAddressIdentifier
	default:
		return a.GetAlias(mainAddressIdentifier), mainAddressIdentifier
	}
}

func (a *Account) RequestAliasAddress(mainAddress []byte, mainAddressIdentifier core.AddressIdentifier, request *vmcommon.AddressRequest) ([]byte, error) {
	aliasAddress := a.GetAlias(request.RequestedIdentifier)
	if aliasAddress != nil {
		return aliasAddress, nil
	}

	requestedAddress, err := address.GeneratePseudoAddress(mainAddress, mainAddressIdentifier, request.RequestedIdentifier)
	if err != nil {
		return nil, err
	}

	if request.SaveOnGenerate {
		a.SetAlias(requestedAddress, request.RequestedIdentifier)
	}
	return requestedAddress, nil
}
