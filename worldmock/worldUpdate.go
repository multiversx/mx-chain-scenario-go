package worldmock

import (
	"errors"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// UpdateBalance sets a new balance to an account
func (b *MockWorld) UpdateBalance(address []byte, newBalance *big.Int) error {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return errors.New("method UpdateBalance expects an existing address")
	}
	acct.Balance = newBalance
	return nil
}

// UpdateBalanceWithDelta changes balance of an account by a given amount
func (b *MockWorld) UpdateBalanceWithDelta(address []byte, balanceDelta *big.Int) error {
	acct := b.AcctMap.GetAccount(address)
	if acct == nil {
		return errors.New("method UpdateBalanceWithDelta expects an existing address")
	}
	acct.Balance = big.NewInt(0).Add(acct.Balance, balanceDelta)
	return nil
}

// UpdateWorldStateBefore performs gas payment, before transaction
func (b *MockWorld) UpdateWorldStateBefore(
	fromAddr []byte,
	gasLimit uint64,
	gasPrice uint64) error {

	acct := b.AcctMap.GetAccount(fromAddr)
	if acct == nil {
		return errors.New("method UpdateWorldStateBefore expects an existing address")
	}
	acct.Nonce++
	gasPayment := big.NewInt(0).Mul(
		big.NewInt(0).SetUint64(gasLimit),
		big.NewInt(0).SetUint64(gasPrice))
	if acct.Balance.Cmp(gasPayment) < 0 {
		return errors.New("not enough balance to pay gas upfront")
	}
	acct.Balance.Sub(acct.Balance, gasPayment)
	return nil
}

// UpdateAccounts should be called after the VM test has run, to update world state
func (b *MockWorld) UpdateAccounts(
	outputAccounts map[string]vmcommon.OutputAccountHandler,
	accountsToDelete [][]byte) error {

	for _, modAcct := range outputAccounts {
		b.UpdateAccountFromOutputAccount(modAcct)
	}

	for _, delAddr := range accountsToDelete {
		b.AcctMap.DeleteAccount(delAddr)
	}

	return nil
}

// UpdateAccountFromOutputAccount updates a single account from a transaction output.
func (b *MockWorld) UpdateAccountFromOutputAccount(modAcct vmcommon.OutputAccountHandler) {
	acct := b.AcctMap.GetAccount(modAcct.GetAddress())
	if acct == nil {
		acct = b.AcctMap.CreateAccount(modAcct.GetAddress(), b)
		acct.OwnerAddress = modAcct.GetCodeDeployerAddress()
		b.AcctMap.PutAccount(acct)
	}
	acct.Exists = true
	if modAcct.GetBalanceDelta() != nil {
		acct.Balance = big.NewInt(0).Add(acct.Balance, modAcct.GetBalanceDelta())
	} else {
		acct.Balance = modAcct.GetBalance()
	}
	if modAcct.GetNonce() > acct.Nonce {
		acct.Nonce = modAcct.GetNonce()
	}
	if len(modAcct.GetCode()) > 0 {
		codeMetadata := vmcommon.CodeMetadataFromBytes(modAcct.GetCodeMetadata())
		acct.SetCodeAndMetadata(modAcct.GetCode(), &codeMetadata)
	}
	if len(modAcct.GetOutputTransfers()) > 0 && len(modAcct.GetOutputTransfers()[0].Data) > 0 {
		acct.AsyncCallData = string(modAcct.GetOutputTransfers()[0].Data)
	}

	for _, stu := range modAcct.GetStorageUpdates() {
		acct.Storage[string(stu.Offset)] = stu.Data
	}
}

// CreateStateBackup -
func (b *MockWorld) CreateStateBackup() {
	b.AccountsAdapter.(*MockAccountsAdapter).SnapshotState(nil, nil)
}

// CommitChanges -
func (b *MockWorld) CommitChanges() error {
	_, err := b.AccountsAdapter.Commit()
	return err
}

// RollbackChanges should be called after the VM test has run, if the tx has failed
func (b *MockWorld) RollbackChanges() error {
	return b.AccountsAdapter.RevertToSnapshot(0)
}
