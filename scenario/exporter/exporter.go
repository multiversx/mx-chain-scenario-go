package exporter

import (
	"errors"
	"math/big"

	mc "github.com/multiversx/mx-chain-scenario-go/scenario/io"
	mj "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	"github.com/multiversx/mx-chain-scenario-go/worldmock/esdtconvert"
)

var errFirstStepMustSetState = errors.New("first step must be of type SetState")

var errNoStepsProvided = errors.New("no steps were provided")

var errScAccountMustHaveOwner = errors.New("scAccount must have owner")

var okStatus = big.NewInt(0)

// ScAddressPrefix is the smart contract address prefix
var ScAddressPrefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 5, 0}

// ScAddressPrefixLength defines the length of the smart contract address prefix
var ScAddressPrefixLength = 10

var benchmarkTxIdent = "benchmark"

var InvalidBenchmarkTxPos = -1

var minimumAcceptedGasPrice = uint64(1)

// ScenarioWithBenchmark defines the component used to hold scenario benchmark
type ScenarioWithBenchmark struct {
	Accs           []*TestAccount
	DeployedAccs   []*TestAccount
	Txs            []*Transaction
	DeployTxs      []*Transaction
	BenchmarkTxPos int
}

func getInvalidScenarioWithBenchmark() ScenarioWithBenchmark {
	return ScenarioWithBenchmark{
		Accs:           nil,
		DeployedAccs:   nil,
		Txs:            nil,
		DeployTxs:      nil,
		BenchmarkTxPos: InvalidBenchmarkTxPos,
	}
}

// GetAccountsAndTransactionsFromScenarios will retrieve the ScenarioWithBenchmark component
func GetAccountsAndTransactionsFromScenarios(testPath string) (stateAndBenchmarkInfo ScenarioWithBenchmark, err error) {
	scenario, err := getScenario(testPath)
	if err != nil {
		return getInvalidScenarioWithBenchmark(), err
	}
	steps := scenario.Steps
	stateAndBenchmarkInfo, err = getAccountsAndTransactionsFromSteps(steps)
	if err != nil {
		return getInvalidScenarioWithBenchmark(), err
	}
	return stateAndBenchmarkInfo, nil
}

func getScenario(testPath string) (scenario *mj.Scenario, err error) {
	scenario, err = mc.ParseScenariosScenarioDefaultParser(testPath)
	if err != nil {
		return nil, err
	}
	return scenario, err
}

func getAccountsAndTransactionsFromSteps(steps []mj.Step) (stateAndBenchmarkInfo ScenarioWithBenchmark, err error) {
	stateAndBenchmarkInfo.BenchmarkTxPos = InvalidBenchmarkTxPos

	if len(steps) == 0 {
		return getInvalidScenarioWithBenchmark(), errNoStepsProvided
	}
	if !stepIsSetState(steps[0]) && !stepIsExternalStep(steps[0]) {
		return getInvalidScenarioWithBenchmark(), errFirstStepMustSetState
	}

	stateAndBenchmarkInfo.Txs = make([]*Transaction, 0)
	stateAndBenchmarkInfo.DeployTxs = make([]*Transaction, 0)
	stateAndBenchmarkInfo.Accs = make([]*TestAccount, 0)
	stateAndBenchmarkInfo.DeployedAccs = make([]*TestAccount, 0)

	for i := 0; i < len(steps); i++ {
		switch step := steps[i].(type) {
		case *mj.SetStateStep:
			setStepAccounts, setStepDeployedAccounts, err := getAccountsFromSetStateStep(step)
			if err != nil {
				return getInvalidScenarioWithBenchmark(), err
			}
			stateAndBenchmarkInfo.Accs = append(stateAndBenchmarkInfo.Accs, setStepAccounts...)
			stateAndBenchmarkInfo.DeployedAccs = append(stateAndBenchmarkInfo.DeployedAccs, setStepDeployedAccounts...)

		case *mj.TxStep:
			if step.ExpectedResult.Status.Value.Cmp(okStatus) == 0 {

				if step.Tx.GasPrice.Value == 0 {
					step.Tx.GasPrice.Value = minimumAcceptedGasPrice
				}
				arguments := getArguments(step.Tx.Arguments)
				switch step.StepTypeName() {
				case "scCall":
					if txIdRequiresBenchmark(step.TxIdent) && benchmarkTxPosIsNotSet(stateAndBenchmarkInfo.BenchmarkTxPos) {
						stateAndBenchmarkInfo.BenchmarkTxPos = len(stateAndBenchmarkInfo.Txs)
					}
					tx := CreateTransaction(
						step.Tx.Function,
						arguments,
						step.Tx.Nonce.Value,
						step.Tx.EGLDValue.Value,
						step.Tx.ESDTValue,
						step.Tx.From.Value,
						append(ScAddressPrefix, step.Tx.To.Value[ScAddressPrefixLength:]...),
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					stateAndBenchmarkInfo.Txs = append(stateAndBenchmarkInfo.Txs, tx)
				case "scUpgrade":
					if txIdRequiresBenchmark(step.TxIdent) && benchmarkTxPosIsNotSet(stateAndBenchmarkInfo.BenchmarkTxPos) {
						stateAndBenchmarkInfo.BenchmarkTxPos = len(stateAndBenchmarkInfo.Txs)
					}
					tx := CreateUpgradeTransaction(
						arguments,
						step.Tx.Code.Original,
						step.Tx.From.Value,
						append(ScAddressPrefix, step.Tx.To.Value[ScAddressPrefixLength:]...),
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					stateAndBenchmarkInfo.Txs = append(stateAndBenchmarkInfo.Txs, tx)
				case "scDeploy":
					deployTx := CreateDeployTransaction(
						arguments,
						step.Tx.Code.Original,
						step.Tx.From.Value,
						step.Tx.GasLimit.Value,
						step.Tx.GasPrice.Value,
					)
					stateAndBenchmarkInfo.DeployTxs = append(stateAndBenchmarkInfo.DeployTxs, deployTx)
				default:
					steps = append(steps[:i], steps[i+1:]...)
					i--
				}
			}
		case *mj.ExternalStepsStep:
			externalStateAndBenchmarkInfo, err := GetAccountsAndTransactionsFromScenarios(step.Path)
			if err != nil {
				return getInvalidScenarioWithBenchmark(), err
			}
			if benchmarkTxPosIsNotSet(stateAndBenchmarkInfo.BenchmarkTxPos) {
				stateAndBenchmarkInfo.BenchmarkTxPos = externalStateAndBenchmarkInfo.BenchmarkTxPos
			}
			stateAndBenchmarkInfo.Accs = append(stateAndBenchmarkInfo.Accs, externalStateAndBenchmarkInfo.Accs...)
			stateAndBenchmarkInfo.DeployedAccs = append(stateAndBenchmarkInfo.DeployedAccs, externalStateAndBenchmarkInfo.DeployedAccs...)
			stateAndBenchmarkInfo.Txs = append(stateAndBenchmarkInfo.Txs, externalStateAndBenchmarkInfo.Txs...)
			stateAndBenchmarkInfo.DeployTxs = append(stateAndBenchmarkInfo.DeployTxs, externalStateAndBenchmarkInfo.DeployTxs...)
		default:
			steps = append(steps[:i], steps[i+1:]...)
			i--
		}
	}
	return stateAndBenchmarkInfo, nil
}

func getAccountsFromSetStateStep(setStateStep *mj.SetStateStep) (accounts []*TestAccount, deployedAccounts []*TestAccount, err error) {
	accounts = make([]*TestAccount, 0)
	deployedAccounts = make([]*TestAccount, 0)
	for _, scenAccount := range setStateStep.Accounts {
		account, err := convertScenariosAccountToTestAccount(scenAccount)
		if err != nil {
			return nil, nil, err
		}
		if len(account.code) > 0 {
			account.address = append(ScAddressPrefix, account.address[ScAddressPrefixLength:]...)
		}
		accounts = append(accounts, account)
	}
	for _, newScenariosAddressMock := range setStateStep.NewAddressMocks {
		scAddress := append(ScAddressPrefix, newScenariosAddressMock.NewAddress.Value[ScAddressPrefixLength:]...)
		ownerAddress := newScenariosAddressMock.CreatorAddress.Value
		account := SetNewAccount(0, scAddress, big.NewInt(0), make(map[string][]byte), make([]byte, 0), ownerAddress)
		deployedAccounts = append(deployedAccounts, account)
	}

	return accounts, deployedAccounts, nil
}

func convertScenariosAccountToTestAccount(scenAcc *mj.Account) (*TestAccount, error) {
	if len(scenAcc.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}
	storage := make(map[string][]byte)
	for _, stkvp := range scenAcc.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}
	_ = esdtconvert.WriteScenariosESDTToStorage(scenAcc.ESDTData, storage)
	account := SetNewAccount(scenAcc.Nonce.Value, scenAcc.Address.Value, scenAcc.Balance.Value, storage, scenAcc.Code.Value, scenAcc.Owner.Value)

	if len(account.code) != 0 && len(account.ownerAddress) == 0 {
		return nil, errScAccountMustHaveOwner
	}

	return account, nil
}

func getArguments(args []mj.JSONBytesFromTree) [][]byte {
	arguments := make([][]byte, len(args))
	for i := 0; i < len(args); i++ {
		arguments[i] = append(arguments[i], args[i].Value...)
	}
	return arguments
}

func stepIsSetState(step mj.Step) bool {
	return step.StepTypeName() == "setState"
}

func stepIsExternalStep(step mj.Step) bool {
	return step.StepTypeName() == "externalSteps"
}

func benchmarkTxPosIsNotSet(benchmarkTxPos int) bool {
	return benchmarkTxPos == InvalidBenchmarkTxPos
}

func txIdRequiresBenchmark(txIdent string) bool {
	return txIdent == benchmarkTxIdent
}
