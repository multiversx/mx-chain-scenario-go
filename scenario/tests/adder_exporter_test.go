package scenTests

import (
	"math/big"
	"testing"

	exporter "github.com/multiversx/mx-chain-scenario-go/scenario/exporter"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	"github.com/multiversx/mx-chain-scenario-go/util"

	"github.com/stretchr/testify/require"
)

// address:owner
var addressOwner = []byte("owner___________________________")

// address:adder
var addressAdder = []byte("adder___________________________")

// address:alice
var addressAlice = []byte("alice___________________________")

// address:bob
var addressBob = []byte("bob_____________________________")

// sc:deployedAdder
var addressDeployedAdder = []byte("deployedAdder___________________")

func TestGetAccountsAndTransactionsFrom_Adder(t *testing.T) {
	sbi, err := exporter.GetAccountsAndTransactionsFromScenarios("adder.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*exporter.TestAccount, 0)
	expectedDeployedAccs := make([]*exporter.TestAccount, 0)
	expectedTxs := make([]*exporter.Transaction, 0)
	expectedDeployTxs := make([]*exporter.Transaction, 0)
	expectedBenchmarkTxPos := 1

	ownerAccount := exporter.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := exporter.SetNewAccount(0, append(exporter.ScAddressPrefix, addressAdder[exporter.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), util.GetSCCode("adder.wasm"), addressOwner)
	deployedScAccount := exporter.SetNewAccount(0, append(exporter.ScAddressPrefix, addressDeployedAdder[exporter.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), make([]byte, 0), addressOwner)
	expectedAccs = append(expectedAccs, ownerAccount, scAccount)
	expectedDeployedAccs = append(expectedDeployedAccs, deployedScAccount)

	transaction := exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*scenmodel.ESDTTxData, 0), sbi.Accs[0].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	expectedTxs = append(expectedTxs, transaction, transaction)

	require.Nil(t, err)
	require.Equal(t, expectedBenchmarkTxPos, sbi.BenchmarkTxPos)
	require.Equal(t, expectedAccs, sbi.Accs)
	require.Equal(t, expectedDeployedAccs, sbi.DeployedAccs)
	require.Equal(t, expectedDeployTxs, sbi.DeployTxs)
	require.Equal(t, expectedTxs, sbi.Txs)
}

func TestGetAccountsAndTransactionsFrom_AdderWithExternalSteps(t *testing.T) {
	sbi, err := exporter.GetAccountsAndTransactionsFromScenarios("adder_with_external_steps.scen.json")
	require.Nil(t, err)
	expectedAccs := make([]*exporter.TestAccount, 0)
	expectedTxs := make([]*exporter.Transaction, 0)
	expectedDeployTxs := make([]*exporter.Transaction, 0)
	expectedBenchmarkTxPos := 1

	ownerAccount := exporter.SetNewAccount(1, addressOwner, big.NewInt(48), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	scAccount := exporter.SetNewAccount(0, append(exporter.ScAddressPrefix, addressAdder[exporter.ScAddressPrefixLength:]...), big.NewInt(0), make(map[string][]byte), util.GetSCCode("adder.wasm"), addressOwner)
	aliceAccount := exporter.SetNewAccount(5, addressAlice, big.NewInt(284), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	bobAccount := exporter.SetNewAccount(3, addressBob, big.NewInt(11), make(map[string][]byte), make([]byte, 0), make([]byte, 0))
	expectedAccs = append(expectedAccs, aliceAccount, scAccount, bobAccount, ownerAccount)
	require.Equal(t, expectedAccs, sbi.Accs)

	transactionAlice := exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*scenmodel.ESDTTxData, 0), sbi.Accs[0].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	transactionBob := exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*scenmodel.ESDTTxData, 0), sbi.Accs[2].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	transactionOwner := exporter.CreateTransaction("add", [][]byte{{3}}, 0, big.NewInt(0), make([]*scenmodel.ESDTTxData, 0), sbi.Accs[3].GetAddress(), sbi.Accs[1].GetAddress(), 5000000, 1)
	expectedTxs = append(expectedTxs, transactionBob, transactionAlice, transactionOwner)
	require.Equal(t, expectedBenchmarkTxPos, sbi.BenchmarkTxPos)
	require.Equal(t, expectedTxs, sbi.Txs)
	require.Equal(t, expectedDeployTxs, sbi.DeployTxs)
}
