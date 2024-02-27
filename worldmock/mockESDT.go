package worldmock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-scenario-go/worldmock/esdtconvert"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// GetTokenBalance returns the ESDT balance of an account for the given token
// key (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenBalance(address []byte, tokenIdentifier []byte, nonce uint64) (*big.Int, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return big.NewInt(0), nil
	}
	return esdtconvert.GetTokenBalance(tokenIdentifier, nonce, account.Storage)
}

// GetTokenData gets the ESDT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) GetTokenData(address []byte, tokenIdentifier []byte, nonce uint64) (*esdt.ESDigitalToken, error) {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return &esdt.ESDigitalToken{
			Value: big.NewInt(0),
		}, nil
	}
	systemAccStorage := make(map[string][]byte)
	systemAcc := bf.World.AcctMap.GetAccount(vmcommon.SystemAccountAddress)
	if systemAcc != nil {
		systemAccStorage = systemAcc.Storage
	}
	return account.GetTokenData(tokenIdentifier, nonce, systemAccStorage)
}

// SetTokenData sets the ESDT information related to a token from the storage of an account
// (token keys are built from the token identifier using MakeTokenKey).
func (bf *BuiltinFunctionsWrapper) SetTokenData(address []byte, tokenIdentifier []byte, nonce uint64, tokenData *esdt.ESDigitalToken) error {
	account := bf.World.AcctMap.GetAccount(address)
	if check.IfNil(account) {
		return nil
	}
	return account.SetTokenData(tokenIdentifier, nonce, tokenData)
}

// ConvertToBuiltinFunction converts a VM input with a populated ESDT field into a built-in function call.
func ConvertToBuiltinFunction(tx *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	switch len(tx.ESDTTransfers) {
	case 0:
		return tx
	case 1:
		return convertToESDTTransfer(tx, tx.ESDTTransfers[0])
	default:
		return convertToMultiESDTTransfer(tx)
	}
}

// PerformDirectESDTTransfer calls the real ESDTTransfer function immediately;
// only works for in-shard transfers for now, but it will be expanded to
// cross-shard.
func convertToESDTTransfer(tx *vmcommon.ContractCallInput, esdtTransfer *vmcommon.ESDTTransfer) *vmcommon.ContractCallInput {
	esdtTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  tx.CallerAddr,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    tx.CallType,
			GasPrice:    tx.GasPrice,
			GasProvided: tx.GasProvided,
			GasLocked:   tx.GasLocked,
		},
		RecipientAddr:     tx.RecipientAddr,
		Function:          core.BuiltInFunctionESDTTransfer,
		AllowInitFunction: false,
	}

	if esdtTransfer.ESDTTokenNonce > 0 {
		esdtTransferInput.Function = core.BuiltInFunctionESDTNFTTransfer
		esdtTransferInput.RecipientAddr = esdtTransferInput.CallerAddr
		nonceAsBytes := big.NewInt(0).SetUint64(esdtTransfer.ESDTTokenNonce).Bytes()
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments,
			esdtTransfer.ESDTTokenName, nonceAsBytes, esdtTransfer.ESDTValue.Bytes(), tx.RecipientAddr)
	} else {
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments,
			esdtTransfer.ESDTTokenName, esdtTransfer.ESDTValue.Bytes())
	}

	if len(tx.Function) > 0 {
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, []byte(tx.Function))
		esdtTransferInput.Arguments = append(esdtTransferInput.Arguments, tx.Arguments...)
	}

	return esdtTransferInput
}

func convertToMultiESDTTransfer(tx *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	multiTransferInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  tx.CallerAddr,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    tx.CallType,
			GasPrice:    tx.GasPrice,
			GasProvided: tx.GasProvided,
			GasLocked:   tx.GasLocked,
		},
		RecipientAddr:     tx.CallerAddr,
		Function:          core.BuiltInFunctionMultiESDTNFTTransfer,
		AllowInitFunction: false,
	}

	multiTransferInput.Arguments = append(multiTransferInput.Arguments, tx.RecipientAddr)

	nrTransfers := len(tx.ESDTTransfers)
	nrTransfersAsBytes := big.NewInt(0).SetUint64(uint64(nrTransfers)).Bytes()
	multiTransferInput.Arguments = append(multiTransferInput.Arguments, nrTransfersAsBytes)

	for _, esdtTransfer := range tx.ESDTTransfers {
		token := esdtTransfer.ESDTTokenName
		nonceAsBytes := big.NewInt(0).SetUint64(esdtTransfer.ESDTTokenNonce).Bytes()
		value := esdtTransfer.ESDTValue

		multiTransferInput.Arguments = append(multiTransferInput.Arguments, token, nonceAsBytes, value.Bytes())
	}

	if len(tx.Function) > 0 {
		multiTransferInput.Arguments = append(multiTransferInput.Arguments, []byte(tx.Function))
		multiTransferInput.Arguments = append(multiTransferInput.Arguments, tx.Arguments...)
	}

	return multiTransferInput
}
