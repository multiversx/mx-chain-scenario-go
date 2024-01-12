package scenexec

import (
	"errors"

	mj "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// ExecuteSetStateStep executes a SetStateStep.
func (ae *ScenarioExecutor) ExecuteSetStateStep(step *mj.SetStateStep) error {
	if len(step.Comment) > 0 {
		log.Trace("SetStateStep", "comment", step.Comment)
	}

	for _, scenAccount := range step.Accounts {
		if scenAccount.Update {
			err := ae.UpdateAccount(scenAccount)
			if err != nil {
				log.Debug("could not update account", err)
				return err
			}
		} else {
			err := ae.PutNewAccount(scenAccount)
			if err != nil {
				log.Debug("could not put new account", err)
				return err
			}
		}
	}

	// replace block info
	ae.World.PreviousBlockInfo = convertBlockInfo(step.PreviousBlockInfo, ae.World.PreviousBlockInfo)
	ae.World.CurrentBlockInfo = convertBlockInfo(step.CurrentBlockInfo, ae.World.CurrentBlockInfo)
	ae.World.Blockhashes = step.BlockHashes.ToValues()

	// append NewAddressMocks
	err := validateNewAddressMocks(step.NewAddressMocks)
	if err != nil {
		return err
	}
	addressMocksToAdd := convertNewAddressMocks(step.NewAddressMocks)
	ae.World.NewAddressMocks = append(ae.World.NewAddressMocks, addressMocksToAdd...)

	return nil
}

// PutNewAccount Puts a new account in world account map. Overwrites.
func (ae *ScenarioExecutor) PutNewAccount(scenAccount *mj.Account) error {
	worldAccount, err := convertAccount(scenAccount, ae.World)
	if err != nil {
		return err
	}
	err = validateSetStateAccount(scenAccount, worldAccount)
	if err != nil {
		return err
	}

	ae.World.AcctMap.PutAccount(worldAccount)
	return nil
}

// UpdateAccount Updates an account in world account map.
func (ae *ScenarioExecutor) UpdateAccount(scenAccount *mj.Account) error {
	worldAccount, err := convertAccount(scenAccount, ae.World)
	if err != nil {
		return err
	}
	err = validateSetStateAccount(scenAccount, worldAccount)
	if err != nil {
		return err
	}

	existingAccount := ae.World.AcctMap.GetAccount(scenAccount.Address.Value)
	if existingAccount == nil {
		return errors.New("account not found. could not update")
	}

	for k, v := range worldAccount.Storage {
		existingAccount.Storage[k] = v
	}
	if !scenAccount.Nonce.Unspecified {
		existingAccount.Nonce = worldAccount.Nonce
	}
	if !scenAccount.Balance.Unspecified {
		existingAccount.Balance = worldAccount.Balance
	}
	if !scenAccount.Username.Unspecified {
		existingAccount.Username = worldAccount.Username
	}
	if !scenAccount.Owner.Unspecified {
		existingAccount.OwnerAddress = worldAccount.OwnerAddress
	}
	if !scenAccount.Code.Unspecified {
		existingAccount.Code = worldAccount.Code
	}
	if !scenAccount.Shard.Unspecified {
		existingAccount.ShardID = worldAccount.ShardID
	}
	existingAccount.AsyncCallData = worldAccount.AsyncCallData

	ae.World.AcctMap.PutAccount(existingAccount)
	return nil
}
