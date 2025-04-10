package scenexec

import (
	"errors"
	"fmt"
	"github.com/multiversx/mx-chain-core-go/core"
	"math/big"

	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-scenario-go/worldmock/esdtconvert"
)

// DefaultCodeMetadata indicates what code metadata to use in smart contracts if unspecified
var DefaultCodeMetadata = []byte{0x05, 0x06}

// ExecuteSetStateStep executes a SetStateStep.
func (ae *ScenarioExecutor) ExecuteSetStateStep(step *scenmodel.SetStateStep) error {
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
func (ae *ScenarioExecutor) PutNewAccount(scenAccount *scenmodel.Account) error {
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
func (ae *ScenarioExecutor) UpdateAccount(scenAccount *scenmodel.Account) error {
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

	codeMetadata := testAcct.CodeMetadata.Value
	if len(testAcct.Code.Value) > 0 && testAcct.CodeMetadata.Unspecified {
		codeMetadata = DefaultCodeMetadata
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
		CodeMetadata:    codeMetadata,
		OwnerAddress:    testAcct.Owner.Value,
		AsyncCallData:   testAcct.AsyncCallData,
		ShardID:         uint32(testAcct.Shard.Value),
		IsSmartContract: len(testAcct.Code.Value) > 0,
		Aliases:         make(map[core.AddressIdentifier][]byte),
		MockWorld:       world,
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
