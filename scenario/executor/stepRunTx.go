package scenexec

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// ExecuteTxStep executes a TxStep.
func (ae *ScenarioExecutor) ExecuteTxStep(step *scenmodel.TxStep) (*vmcommon.VMOutput, error) {
	log.Trace("ExecuteTxStep", "id", step.TxIdent)
	if len(step.Comment) > 0 {
		log.Trace("ExecuteTxStep", "comment", step.Comment)
	}

	if step.DisplayLogs {
		SetLoggingForTests()
	}

	output, err := ae.executeTx(step.TxIdent, step.Tx)
	if err != nil {
		return nil, err
	}

	if step.DisplayLogs {
		DisableLoggingForTests()
	}

	// check results
	if step.ExpectedResult != nil {
		err = ae.checkTxResults(step.TxIdent, step.ExpectedResult, ae.checkGas, output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func (ae *ScenarioExecutor) executeTx(txIndex string, tx *scenmodel.Transaction) (*vmcommon.VMOutput, error) {
	var err error
	gasForExecution := uint64(0)

	// use gas (before snaphot)
	if tx.Type.HasSender() {
		beforeErr := ae.World.UpdateWorldStateBefore(
			tx.From.Value,
			tx.GasLimit.Value,
			tx.GasPrice.Value)
		if beforeErr != nil {
			err = fmt.Errorf("could not set up tx %s: %w", txIndex, beforeErr)
			return nil, err
		}

		gasForExecution = tx.GasLimit.Value
	}

	ae.World.CreateStateBackup()

	defer func() {
		if err != nil {
			errRollback := ae.World.RollbackChanges()
			if errRollback != nil {
				err = errRollback
			}
		} else {
			errCommit := ae.World.CommitChanges()
			if errCommit != nil {
				err = errCommit
			}
		}
	}()

	// we also use fake vm outputs for transactions that don't use the VM, just for convenience
	var output *vmcommon.VMOutput

	if !ae.senderHasEnoughBalance(tx) {
		// out of funds is handled by the protocol, so it needs to be mocked here
		output = outOfFundsResult()
	} else {
		switch tx.Type {
		case scenmodel.ScDeploy:
			output, err = ae.scCreate(txIndex, tx, gasForExecution)
			if err != nil {
				return nil, err
			}
			if ae.PeekTraceGas() {
				fmt.Println("\nIn txID:", txIndex, ", step type:Deploy", ", total gas used:", gasForExecution-output.GasRemaining)
			}
		case scenmodel.ScQuery:
			// imitates the behaviour of the protocol
			// the sender is the contract itself during SC queries
			tx.From = tx.To
			// gas restrictions waived during SC queries
			tx.GasLimit.Value = math.MaxInt64
			gasForExecution = math.MaxInt64
			fallthrough
		case scenmodel.ScCall:
			output, err = ae.scCall(txIndex, tx, tx.GasLimit.Value)
			if err != nil {
				return nil, err
			}
			if ae.PeekTraceGas() {
				fmt.Println("\nIn txID:", txIndex, ", step type:ScCall, function:", tx.Function, ", total gas used:", gasForExecution-output.GasRemaining)
			}
		case scenmodel.Transfer:
			if tx.ESDTValue != nil {
				output, err = ae.directESDTTransfer(tx)
				if err != nil {
					return nil, err
				}
			} else {
				output = ae.simpleTransferOutput(tx)
			}
		case scenmodel.ValidatorReward:
			output, err = ae.validatorRewardOutput(tx)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown transaction type")
		}
	}

	if output.ReturnCode == vmcommon.Ok {
		err := ae.updateStateAfterTx(tx, output)
		if err != nil {
			return nil, err
		}
	} else {
		err = fmt.Errorf(
			"tx step failed: retcode=%d, msg=%s",
			output.ReturnCode, output.ReturnMessage)
	}

	return output, nil
}

func (ae *ScenarioExecutor) senderHasEnoughBalance(tx *scenmodel.Transaction) bool {
	if !tx.Type.HasSender() {
		return true
	}
	sender := ae.World.AcctMap.GetAccount(tx.From.Value)
	return sender.Balance.Cmp(tx.EGLDValue.Value) >= 0
}

func (ae *ScenarioExecutor) simpleTransferOutput(tx *scenmodel.Transaction) *vmcommon.VMOutput {
	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(tx.To.Value)] = &vmcommon.OutputAccount{
		Address:      tx.To.Value,
		BalanceDelta: tx.EGLDValue.Value,
	}

	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  outputAccounts,
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func (ae *ScenarioExecutor) validatorRewardOutput(tx *scenmodel.Transaction) (*vmcommon.VMOutput, error) {
	reward := tx.EGLDValue.Value
	recipient := ae.World.AcctMap.GetAccount(tx.To.Value)
	if recipient == nil {
		return nil, fmt.Errorf("tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To.Value))
	}
	recipient.BalanceDelta = reward
	storageReward := recipient.StorageValue(RewardKey)
	storageReward = big.NewInt(0).Add(
		big.NewInt(0).SetBytes(storageReward),
		reward).Bytes()

	outputAccounts := make(map[string]*vmcommon.OutputAccount)
	outputAccounts[string(tx.To.Value)] = &vmcommon.OutputAccount{
		Address:      tx.To.Value,
		BalanceDelta: tx.EGLDValue.Value,
		StorageUpdates: map[string]*vmcommon.StorageUpdate{
			RewardKey: {
				Offset: []byte(RewardKey),
				Data:   storageReward,
			},
		},
	}

	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  outputAccounts,
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}, nil
}

func outOfFundsResult() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.OutOfFunds,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func (ae *ScenarioExecutor) scCreate(txIndex string, tx *scenmodel.Transaction, gasLimit uint64) (*vmcommon.VMOutput, error) {
	txHash := generateTxHash(txIndex)
	vmInput := vmcommon.VMInput{
		CallerAddr:     tx.From.Value,
		Arguments:      scenmodel.JSONBytesFromTreeValues(tx.Arguments),
		CallValue:      tx.EGLDValue.Value,
		GasPrice:       tx.GasPrice.Value,
		GasProvided:    gasLimit,
		OriginalTxHash: txHash,
		CurrentTxHash:  txHash,
		ESDTTransfers:  make([]*vmcommon.ESDTTransfer, 0),
	}
	addESDTToVMInput(tx.ESDTValue, &vmInput)
	codeMetadata := tx.CodeMetadata.Value
	if tx.CodeMetadata.Unspecified {
		codeMetadata = DefaultCodeMetadata
	}
	input := &vmcommon.ContractCreateInput{
		ContractCode:         tx.Code.Value,
		ContractCodeMetadata: codeMetadata,
		VMInput:              vmInput,
	}

	return ae.vm.RunSmartContractCreate(input)
}

func (ae *ScenarioExecutor) scCall(txIndex string, tx *scenmodel.Transaction, gasLimit uint64) (*vmcommon.VMOutput, error) {
	recipient := ae.World.AcctMap.GetAccount(tx.To.Value)
	if recipient == nil {
		return nil, fmt.Errorf("tx recipient (address: %s) does not exist", hex.EncodeToString(tx.To.Value))
	}
	if len(recipient.Code) == 0 {
		return nil, fmt.Errorf("tx recipient (address: %s) is not a smart contract", hex.EncodeToString(tx.To.Value))
	}

	input := ConvertScenarioTxToVMInput(tx)
	txHash := generateTxHash(txIndex)
	input.CurrentTxHash = txHash
	input.OriginalTxHash = txHash

	if len(input.ESDTTransfers) > 0 {
		bfInput := worldmock.ConvertToBuiltinFunction(input)
		vmOutput, err := ae.World.BuiltinFuncs.ProcessBuiltInFunction(bfInput)
		if err != nil {
			return nil, err
		}
		if vmOutput.ReturnCode != vmcommon.Ok {
			return nil, fmt.Errorf(
				"%s failed: retcode = %d, msg = %s",
				bfInput.Function,
				vmOutput.ReturnCode,
				vmOutput.ReturnMessage)
		}
	}

	return ae.vm.RunSmartContractCall(input)
}

func (ae *ScenarioExecutor) directESDTTransfer(tx *scenmodel.Transaction) (*vmcommon.VMOutput, error) {
	input := ConvertScenarioTxToVMInput(tx)
	bfInput := worldmock.ConvertToBuiltinFunction(input)
	vmOutput, err := ae.World.BuiltinFuncs.ProcessBuiltInFunction(bfInput)

	if err != nil {
		return nil, err
	}

	if vmOutput.ReturnCode != vmcommon.Ok {
		return nil, fmt.Errorf(
			"%s failed: retcode = %d, msg = %s",
			bfInput.Function,
			vmOutput.ReturnCode,
			vmOutput.ReturnMessage)
	}

	return vmOutput, err
}

func (ae *ScenarioExecutor) updateStateAfterTx(
	tx *scenmodel.Transaction,
	output *vmcommon.VMOutput) error {

	// subtract call value from sender (this is not reflected in the delta)
	// except for validatorReward, there is no sender there
	if tx.Type.HasSender() {
		_ = ae.World.UpdateBalanceWithDelta(tx.From.Value, big.NewInt(0).Neg(tx.EGLDValue.Value))
	}

	// update accounts based on deltas
	updErr := ae.World.UpdateAccounts(output.OutputAccounts, output.DeletedAccounts)
	if updErr != nil {
		return updErr
	}

	// sum of all balance deltas should equal call value (unless we got an error)
	// (unless it is validatorReward, when funds just pop into existence)
	if tx.Type.HasSender() {
		sumOfBalanceDeltas := big.NewInt(0)
		for _, oa := range output.OutputAccounts {
			sumOfBalanceDeltas = sumOfBalanceDeltas.Add(sumOfBalanceDeltas, oa.BalanceDelta)
		}
		if sumOfBalanceDeltas.Cmp(tx.EGLDValue.Value) != 0 {
			return fmt.Errorf("sum of balance deltas should equal call value. Sum of balance deltas: %d (0x%x). Call value: %d (0x%x)",
				sumOfBalanceDeltas, sumOfBalanceDeltas, tx.EGLDValue.Value, tx.EGLDValue.Value)
		}
	}

	return nil
}

func generateTxHash(txIndex string) []byte {
	txIndexBytes := []byte(txIndex)
	if len(txIndexBytes) > 32 {
		return txIndexBytes[:32]
	}
	for i := len(txIndexBytes); i < 32; i++ {
		txIndexBytes = append(txIndexBytes, '.')
	}
	return txIndexBytes
}

func addESDTToVMInput(esdtData []*scenmodel.ESDTTxData, vmInput *vmcommon.VMInput) {
	esdtDataLen := len(esdtData)

	if esdtDataLen > 0 {
		vmInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, esdtDataLen)
		for i := 0; i < esdtDataLen; i++ {
			vmInput.ESDTTransfers[i] = &vmcommon.ESDTTransfer{}
			vmInput.ESDTTransfers[i].ESDTTokenName = esdtData[i].TokenIdentifier.Value
			vmInput.ESDTTransfers[i].ESDTValue = esdtData[i].Value.Value
			vmInput.ESDTTransfers[i].ESDTTokenNonce = esdtData[i].Nonce.Value
			if vmInput.ESDTTransfers[i].ESDTTokenNonce != 0 {
				vmInput.ESDTTransfers[i].ESDTTokenType = uint32(core.NonFungible)
			} else {
				vmInput.ESDTTransfers[i].ESDTTokenType = uint32(core.Fungible)
			}
		}
	}
}

// ConvertScenarioTxToVMInput converts the scenario format to the VM input format.
func ConvertScenarioTxToVMInput(tx *scenmodel.Transaction) *vmcommon.ContractCallInput {
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  tx.From.Value,
			Arguments:   scenmodel.JSONBytesFromTreeValues(tx.Arguments),
			CallValue:   tx.EGLDValue.Value,
			CallType:    vm.DirectCall,
			GasPrice:    tx.GasPrice.Value,
			GasProvided: tx.GasLimit.Value,
			GasLocked:   0,
		},
		RecipientAddr:     tx.To.Value,
		Function:          tx.Function,
		AllowInitFunction: false,
	}
	addESDTToVMInput(tx.ESDTValue, &input.VMInput)
	return input
}
