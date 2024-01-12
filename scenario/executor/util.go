package scenexec

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	er "github.com/multiversx/mx-chain-scenario-go/scenario/expression/reconstructor"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-scenario-go/worldmock/esdtconvert"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

func convertAccount(testAcct *scenmodel.Account, world *worldmock.MockWorld) (*worldmock.Account, error) {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	err := esdtconvert.WriteScenariosESDTToStorage(testAcct.ESDTData, storage)
	if err != nil {
		return nil, err
	}

	if len(testAcct.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}

	account := &worldmock.Account{
		Address:         testAcct.Address.Value,
		Nonce:           testAcct.Nonce.Value,
		Balance:         big.NewInt(0).Set(testAcct.Balance.Value),
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0).Set(testAcct.DeveloperReward.Value),
		Username:        testAcct.Username.Value,
		Storage:         storage,
		Code:            testAcct.Code.Value,
		OwnerAddress:    testAcct.Owner.Value,
		AsyncCallData:   testAcct.AsyncCallData,
		ShardID:         uint32(testAcct.Shard.Value),
		IsSmartContract: len(testAcct.Code.Value) > 0,
		CodeMetadata: (&vmcommon.CodeMetadata{
			Payable:     true,
			Upgradeable: true,
			Readable:    true,
		}).ToBytes(), // TODO: add explicit fields in scenario JSON
		MockWorld: world,
	}

	return account, nil
}

func validateSetStateAccount(scenAccount *scenmodel.Account, converted *worldmock.Account) error {
	err := converted.Validate()
	if err != nil {
		return fmt.Errorf(
			`"setState" step validation failed for account "%s": %w`,
			scenAccount.Address.Original,
			err)
	}
	return nil
}

func validateNewAddressMocks(testNAMs []*scenmodel.NewAddressMock) error {
	for _, testNAM := range testNAMs {
		if !worldmock.IsSmartContractAddress(testNAM.NewAddress.Value) {
			return fmt.Errorf(
				`address in "setState" "newAddresses" field should have SC format: %s`,
				testNAM.NewAddress.Original)
		}
	}
	return nil
}

func convertNewAddressMocks(testNAMs []*scenmodel.NewAddressMock) []*worldmock.NewAddressMock {
	var result []*worldmock.NewAddressMock
	for _, testNAM := range testNAMs {
		result = append(result, &worldmock.NewAddressMock{
			CreatorAddress: testNAM.CreatorAddress.Value,
			CreatorNonce:   testNAM.CreatorNonce.Value,
			NewAddress:     testNAM.NewAddress.Value,
		})
	}
	return result
}

func convertBlockInfo(testBlockInfo *scenmodel.BlockInfo, currentInfo *worldmock.BlockInfo) *worldmock.BlockInfo {
	if testBlockInfo == nil {
		return currentInfo
	}

	if currentInfo == nil {
		currentInfo = &worldmock.BlockInfo{
			BlockTimestamp: 0,
			BlockNonce:     0,
			BlockRound:     0,
			BlockEpoch:     0,
			RandomSeed:     nil,
		}
	}

	if !testBlockInfo.BlockTimestamp.OriginalEmpty() {
		currentInfo.BlockTimestamp = testBlockInfo.BlockTimestamp.Value

	}

	if !testBlockInfo.BlockNonce.OriginalEmpty() {
		currentInfo.BlockNonce = testBlockInfo.BlockNonce.Value
	}

	if !testBlockInfo.BlockRound.OriginalEmpty() {
		currentInfo.BlockRound = testBlockInfo.BlockRound.Value
	}

	if !testBlockInfo.BlockEpoch.OriginalEmpty() {
		currentInfo.BlockEpoch = uint32(testBlockInfo.BlockEpoch.Value)
	}

	if testBlockInfo.BlockRandomSeed != nil && !testBlockInfo.BlockRandomSeed.OriginalEmpty() {
		var randomsSeed [48]byte
		copy(randomsSeed[:], testBlockInfo.BlockRandomSeed.Value)
		currentInfo.RandomSeed = &randomsSeed

	}

	return currentInfo
}

// this is a small hack, so we can reuse JSON printing in error messages
func (ae *ScenarioExecutor) convertLogToTestFormat(outputLog *vmcommon.LogEntry) *scenmodel.LogEntry {
	topics := scenmodel.JSONCheckValueList{
		Values: make([]scenmodel.JSONCheckBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		topics.Values[i] = scenmodel.JSONCheckBytesReconstructed(
			topic,
			ae.exprReconstructor.Reconstruct(topic,
				er.NoHint))
	}

	dataField := scenmodel.JSONCheckValueList{
		Values: make([]scenmodel.JSONCheckBytes, len(outputLog.Data)),
	}
	for i, data := range outputLog.Data {
		dataField.Values[i] = scenmodel.JSONCheckBytesReconstructed(
			data,
			ae.exprReconstructor.Reconstruct(data,
				er.NoHint))
	}
	testLog := scenmodel.LogEntry{
		Address: scenmodel.JSONCheckBytesReconstructed(
			outputLog.Address,
			ae.exprReconstructor.Reconstruct(outputLog.Address,
				er.AddressHint)),
		Endpoint: scenmodel.JSONCheckBytesReconstructed(
			outputLog.Identifier,
			ae.exprReconstructor.Reconstruct(outputLog.Identifier,
				er.StrHint)),
		Data:   dataField,
		Topics: topics,
	}

	return &testLog
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

func logGasTrace(ae *ScenarioExecutor) {
	if ae.PeekTraceGas() {
		scGasTrace := ae.vm.GetGasTrace()
		totalGasUsedByAPIs := 0
		for scAddress, gasTrace := range scGasTrace {
			fmt.Println("Gas Trace for: ", "SC Address", scAddress)
			for functionName, value := range gasTrace {
				totalGasUsed := uint64(0)
				for _, usedGas := range value {
					totalGasUsed += usedGas
				}
				fmt.Println("GasTrace: functionName:", functionName, ",  totalGasUsed:", totalGasUsed, ", numberOfCalls:", len(value))
				totalGasUsedByAPIs += int(totalGasUsed)
			}
			fmt.Println("TotalGasUsedByAPIs: ", totalGasUsedByAPIs)
		}
	}
}

func setGasTraceInMetering(ae *ScenarioExecutor, enable bool) {
	if enable && ae.PeekTraceGas() {
		ae.vm.SetGasTracing(true)
	} else {
		ae.vm.SetGasTracing(false)
	}
}

func setExternalStepGasTracing(ae *ScenarioExecutor, step *scenmodel.ExternalStepsStep) {
	switch step.TraceGas.ToInt() {
	case scenmodel.Undefined.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, ae.PeekTraceGas())
	case scenmodel.TrueValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, true)
	case scenmodel.FalseValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, false)
	}
}

func resetGasTracesIfNewTest(ae *ScenarioExecutor, scenario *scenmodel.Scenario) {
	if ae.vm == nil || scenario.IsNewTest {
		ae.scenarioTraceGas = make([]bool, 0)
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, scenario.TraceGas)
		scenario.IsNewTest = false
	}
}
