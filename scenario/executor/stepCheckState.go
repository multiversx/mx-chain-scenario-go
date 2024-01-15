package scenexec

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	er "github.com/multiversx/mx-chain-scenario-go/scenario/expression/reconstructor"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	worldmock "github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-scenario-go/worldmock/esdtconvert"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// ExecuteCheckStateStep executes a CheckStateStep defined by the current scenario.
func (ae *ScenarioExecutor) ExecuteCheckStateStep(step *scenmodel.CheckStateStep) error {
	if len(step.Comment) > 0 {
		log.Trace("CheckStateStep", "comment", step.Comment)
	}

	baseErrMsg := checkStateBaseErrorMsg(step)
	return ae.checkAccounts(baseErrMsg, step.CheckAccounts)
}

func checkStateBaseErrorMsg(step *scenmodel.CheckStateStep) string {
	if len(step.CheckStateIdent) > 0 {
		return fmt.Sprintf("Check state \"%s\":", step.CheckStateIdent)
	}
	return "Check state:"
}

func (ae *ScenarioExecutor) checkAccounts(baseErrMsg string, checkAccounts *scenmodel.CheckAccounts) error {
	if !checkAccounts.MoreAccountsAllowed {
		for worldAcctAddr := range ae.World.AcctMap {
			postAcctMatch := scenmodel.FindCheckAccount(checkAccounts.Accounts, []byte(worldAcctAddr))
			if postAcctMatch == nil && !bytes.Equal([]byte(worldAcctAddr), vmcommon.SystemAccountAddress) {
				return fmt.Errorf("%s unexpected account address: %s",
					baseErrMsg,
					ae.exprReconstructor.Reconstruct(
						[]byte(worldAcctAddr),
						er.AddressHint))
			}
		}
	}

	for _, expectedAcct := range checkAccounts.Accounts {
		matchingAcct, isMatch := ae.World.AcctMap[string(expectedAcct.Address.Value)]
		if !isMatch {
			return fmt.Errorf("%s account %s expected but not found after running test",
				baseErrMsg,
				expectedAcct.Address.Original)
		}

		if !bytes.Equal(matchingAcct.Address, expectedAcct.Address.Value) {
			return fmt.Errorf("%s bad account address %s",
				baseErrMsg,
				ae.exprReconstructor.Reconstruct(
					matchingAcct.Address,
					er.AddressHint))
		}

		if !expectedAcct.Nonce.Check(matchingAcct.Nonce) {
			return fmt.Errorf("%s bad account nonce. Account: %s. Want: \"%s\". Have: \"%d\"",
				baseErrMsg,
				expectedAcct.Address.Original,
				expectedAcct.Nonce.Original,
				matchingAcct.Nonce)
		}

		if !expectedAcct.Balance.Check(matchingAcct.Balance) {
			return fmt.Errorf("%s bad account balance. Account: %s. Want: \"%s\". Have: \"%s\"",
				baseErrMsg,
				expectedAcct.Address.Original,
				expectedAcct.Balance.Original,
				ae.exprReconstructor.ReconstructFromBigInt(matchingAcct.Balance))
		}

		if !expectedAcct.Username.Check(matchingAcct.Username) {
			return fmt.Errorf("%s bad account username. Account: %s. Want: %s. Have: \"%s\"",
				baseErrMsg,
				expectedAcct.Address.Original,
				oj.JSONString(expectedAcct.Username.Original),
				ae.exprReconstructor.Reconstruct(
					matchingAcct.Username,
					er.StrHint))
		}

		if !expectedAcct.Code.Check(matchingAcct.Code) {
			return fmt.Errorf("%s bad account code. Account: %s. Want: %s. Have: \"%s\"",
				baseErrMsg,
				expectedAcct.Address.Original,
				oj.JSONString(expectedAcct.Code.Original),
				ae.exprReconstructor.Reconstruct(
					matchingAcct.Code,
					er.CodeHint))
		}

		if !expectedAcct.Owner.IsUnspecified() && !bytes.Equal(matchingAcct.OwnerAddress, expectedAcct.Owner.Value) {
			return fmt.Errorf("%s bad account owner. Account: %s. Want: %s. Have: \"%s\"",
				baseErrMsg,
				expectedAcct.Address.Original,
				oj.JSONString(expectedAcct.Owner.Original),
				ae.exprReconstructor.Reconstruct(
					matchingAcct.OwnerAddress,
					er.AddressHint))
		}

		// currently ignoring asyncCallData that is unspecified in the json
		if !expectedAcct.AsyncCallData.IsUnspecified() &&
			!expectedAcct.AsyncCallData.Check([]byte(matchingAcct.AsyncCallData)) {
			return fmt.Errorf("%s bad async call data. Account: %s. Want: [%s]. Have: [%s]",
				baseErrMsg,
				expectedAcct.Address.Original,
				expectedAcct.AsyncCallData.Original,
				matchingAcct.AsyncCallData)
		}

		err := ae.checkAccountStorage(baseErrMsg, expectedAcct, matchingAcct)
		if err != nil {
			return err
		}

		err = ae.checkAccountESDT(baseErrMsg, expectedAcct, matchingAcct)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ae *ScenarioExecutor) checkAccountStorage(baseErrMsg string, expectedAcct *scenmodel.CheckAccount, matchingAcct *worldmock.Account) error {
	if expectedAcct.IgnoreStorage {
		return nil
	}

	expectedStorage := make(map[string]scenmodel.JSONCheckBytes)
	for _, stkvp := range expectedAcct.CheckStorage {
		expectedStorage[string(stkvp.Key.Value)] = stkvp.CheckValue
	}

	allKeys := make(map[string]bool)
	for k := range expectedStorage {
		allKeys[k] = true
	}
	for k := range matchingAcct.Storage {
		allKeys[k] = true
	}
	storageError := ""
	for k := range allKeys {
		// ignore all reserved keys
		if strings.HasPrefix(k, core.ProtectedKeyPrefix) {
			continue
		}

		want, specified := expectedStorage[k]
		if !specified {
			if expectedAcct.MoreStorageAllowed {
				// if `"+": ""` was written in the test, any unspecified entries are allowed,
				// which is equivalent to treating them all as "*".
				want = scenmodel.JSONCheckBytesStar()
			} else {
				// otherwise, by default, any unexpected storage key leads to a test failure
				want = scenmodel.JSONCheckBytesUnspecified()
			}
		}
		have := matchingAcct.StorageValue(k)

		if !want.Check(have) {
			storageError += fmt.Sprintf(
				"\n  for key %s: Want: %s. Have: \"%s\"",
				ae.exprReconstructor.Reconstruct([]byte(k), er.NoHint),
				oj.JSONString(want.Original),
				ae.exprReconstructor.Reconstruct(have, er.NoHint))
		}
	}
	if len(storageError) > 0 {
		return fmt.Errorf("%s wrong account storage for account \"%s\":%s",
			baseErrMsg,
			expectedAcct.Address.Original, storageError)
	}
	return nil
}

func (ae *ScenarioExecutor) checkAccountESDT(baseErrMsg string, expectedAcct *scenmodel.CheckAccount, matchingAcct *worldmock.Account) error {
	if expectedAcct.IgnoreESDT {
		return nil
	}

	systemAccStorage := make(map[string][]byte)
	systemAcc, exists := ae.World.AcctMap[string(vmcommon.SystemAccountAddress)]
	if exists {
		systemAccStorage = systemAcc.Storage
	}

	accountAddress := expectedAcct.Address.Original
	expectedTokens := getExpectedTokens(expectedAcct)
	accountTokens, err := esdtconvert.GetFullMockESDTData(matchingAcct.Storage, systemAccStorage)
	if err != nil {
		return err
	}

	allTokenNames := make(map[string]bool)
	for tokenName := range expectedTokens {
		allTokenNames[tokenName] = true
	}
	for tokenName := range accountTokens {
		allTokenNames[tokenName] = true
	}
	var errs []error
	for tokenName := range allTokenNames {
		expectedToken := expectedTokens[tokenName]
		accountToken := accountTokens[tokenName]
		if expectedToken == nil {
			expectedToken = &scenmodel.CheckESDTData{
				TokenIdentifier: scenmodel.JSONBytesFromString{
					Value:    []byte(tokenName),
					Original: ae.exprReconstructor.Reconstruct([]byte(tokenName), er.StrHint),
				},
				Instances: []*scenmodel.CheckESDTInstance{},
				LastNonce: scenmodel.JSONCheckUint64{Value: 0, Original: ""},
				Roles:     []string{},
			}
		} else if accountToken == nil {
			accountToken = &esdtconvert.MockESDTData{
				TokenIdentifier: []byte(tokenName),
				Instances:       []*esdt.ESDigitalToken{},
				LastNonce:       0,
				Roles:           [][]byte{},
			}
		}

		errs = append(errs, ae.checkTokenState(accountAddress, tokenName, expectedToken, accountToken)...)
	}

	errorString := makeErrorString(errs)
	if len(errorString) > 0 {
		return fmt.Errorf("%s mismatch for account \"%s\":%s", baseErrMsg, accountAddress, errorString)
	}

	return nil
}

func getExpectedTokens(expectedAcct *scenmodel.CheckAccount) map[string]*scenmodel.CheckESDTData {
	expectedTokens := make(map[string]*scenmodel.CheckESDTData)
	for _, expectedTokenData := range expectedAcct.CheckESDTData {
		tokenName := expectedTokenData.TokenIdentifier.Value
		expectedTokens[string(tokenName)] = expectedTokenData
	}

	return expectedTokens
}

func (ae *ScenarioExecutor) checkTokenState(
	accountAddress string,
	tokenName string,
	expectedToken *scenmodel.CheckESDTData,
	accountToken *esdtconvert.MockESDTData,
) []error {

	var errors []error

	errors = append(errors, ae.checkTokenInstances(accountAddress, tokenName, expectedToken, accountToken)...)

	if !expectedToken.LastNonce.Check(accountToken.LastNonce) {
		errors = append(errors, fmt.Errorf("bad account ESDT last nonce. Account: %s. Token: %s. Want: \"%s\". Have: %d",
			accountAddress,
			tokenName,
			expectedToken.LastNonce.Original,
			accountToken.LastNonce))
	}

	errors = append(errors, checkTokenRoles(accountAddress, tokenName, expectedToken, accountToken)...)

	return errors
}

func (ae *ScenarioExecutor) checkTokenInstances(
	_ string,
	tokenName string,
	expectedToken *scenmodel.CheckESDTData,
	accountToken *esdtconvert.MockESDTData,
) []error {

	var errors []error

	allNonces := make(map[uint64]bool)
	expectedInstances := make(map[uint64]*scenmodel.CheckESDTInstance)
	accountInstances := make(map[uint64]*esdt.ESDigitalToken)
	for _, expectedInstance := range expectedToken.Instances {
		nonce := expectedInstance.Nonce.Value
		allNonces[nonce] = true
		expectedInstances[nonce] = expectedInstance
	}
	for _, accountInstance := range accountToken.Instances {
		nonce := accountInstance.TokenMetaData.Nonce
		allNonces[nonce] = true
		accountInstances[nonce] = accountInstance
	}

	for nonce := range allNonces {
		expectedInstance := expectedInstances[nonce]
		accountInstance := accountInstances[nonce]

		if expectedInstance == nil {
			expectedInstance = &scenmodel.CheckESDTInstance{
				Nonce:   scenmodel.JSONUint64{Value: nonce, Original: ""},
				Balance: scenmodel.JSONCheckBigInt{Value: big.NewInt(0), Original: ""},
			}
		} else if accountInstance == nil {
			accountInstance = &esdt.ESDigitalToken{
				Value: big.NewInt(0),
				TokenMetaData: &esdt.MetaData{
					Name:  []byte(tokenName),
					Nonce: nonce,
				},
			}
		}

		if !expectedInstance.Balance.Check(accountInstance.Value) {
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad balance. Want: \"%s\". Have: \"%d\"",
				tokenName,
				nonce,
				expectedInstance.Balance.Original,
				accountInstance.Value))
		}
		if !expectedInstance.Creator.IsUnspecified() &&
			!expectedInstance.Creator.Check(accountInstance.TokenMetaData.Creator) {
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad creator. Want: %s. Have: \"%s\"",
				tokenName,
				nonce,
				objectStringOrDefault(expectedInstance.Creator.Original),
				ae.exprReconstructor.Reconstruct(
					accountInstance.TokenMetaData.Creator,
					er.AddressHint)))
		}
		if !expectedInstance.Royalties.IsUnspecified() &&
			!expectedInstance.Royalties.Check(uint64(accountInstance.TokenMetaData.Royalties)) {
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad royalties. Want: \"%s\". Have: \"%s\"",
				tokenName,
				nonce,
				expectedInstance.Royalties.Original,
				ae.exprReconstructor.ReconstructFromUint64(
					uint64(accountInstance.TokenMetaData.Royalties))))
		}
		if !expectedInstance.Hash.IsUnspecified() &&
			!expectedInstance.Hash.Check(accountInstance.TokenMetaData.Hash) {
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad hash. Want: %s. Have: %s",
				tokenName,
				nonce,
				objectStringOrDefault(expectedInstance.Hash.Original),
				ae.exprReconstructor.Reconstruct(
					accountInstance.TokenMetaData.Hash,
					er.NoHint)))
		}

		if !expectedInstance.Uris.IsUnspecified() &&
			!expectedInstance.Uris.CheckList(accountInstance.TokenMetaData.URIs) {
			// in this case unspecified is interpreted as *
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad URI. Want: %s. Have: %s",
				tokenName,
				nonce,
				checkBytesListPretty(expectedInstance.Uris),
				ae.exprReconstructor.ReconstructList(accountInstance.TokenMetaData.URIs, er.StrHint)))
		}

		if !expectedInstance.Attributes.IsUnspecified() &&
			!expectedInstance.Attributes.Check(accountInstance.TokenMetaData.Attributes) {
			errors = append(errors, fmt.Errorf(
				"for token: %s, nonce: %d: Bad attributes. Want: %s. Have: \"%s\"",
				tokenName,
				nonce,
				objectStringOrDefault(expectedInstance.Attributes.Original),
				ae.exprReconstructor.Reconstruct(
					accountInstance.TokenMetaData.Attributes,
					er.StrHint)))
		}

	}

	return errors
}

func checkTokenRoles(
	accountAddress string,
	tokenName string,
	expectedToken *scenmodel.CheckESDTData,
	accountToken *esdtconvert.MockESDTData) []error {

	var errors []error

	allRoles := make(map[string]bool)
	expectedRoles := make(map[string]bool)
	accountRoles := make(map[string]bool)

	for _, expectedRole := range expectedToken.Roles {
		allRoles[expectedRole] = true
		expectedRoles[expectedRole] = true
	}
	for _, accountRole := range accountToken.Roles {
		allRoles[string(accountRole)] = true
		accountRoles[string(accountRole)] = true
	}
	for role := range allRoles {
		if !expectedRoles[role] {
			errors = append(errors, fmt.Errorf("unexpected ESDT role. Account: %s. Token: %s. Role: %s",
				accountAddress,
				tokenName,
				role))
		}
		if !accountRoles[role] {
			errors = append(errors, fmt.Errorf("missing ESDT role. Account: %s. Token: %s. Role: %s",
				accountAddress,
				tokenName,
				role))
		}
	}

	return errors
}

func makeErrorString(errors []error) string {
	errorString := ""
	for _, err := range errors {
		errorString += "\n  " + err.Error()
	}
	return errorString
}

func objectStringOrDefault(obj oj.OJsonObject) string {
	if obj == nil {
		return ""
	}

	return oj.JSONString(obj)
}
