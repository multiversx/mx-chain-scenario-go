package esdtconvert

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
)

// MakeESDTUserMetadataBytes creates metadata byte slice
func MakeESDTUserMetadataBytes(frozen bool) []byte {
	metadata := &builtInFunctions.ESDTUserMetadata{
		Frozen: frozen,
	}

	return metadata.ToBytes()
}

// WriteScenariosESDTToStorage writes the Scenarios ESDT data to the provided storage map
func WriteScenariosESDTToStorage(esdtData []*scenmodel.ESDTData, destination map[string][]byte) error {
	for _, scenESDTData := range esdtData {
		tokenIdentifier := scenESDTData.TokenIdentifier.Value
		isFrozen := scenESDTData.Frozen.Value > 0
		for _, instance := range scenESDTData.Instances {
			tokenNonce := instance.Nonce.Value
			tokenKey := makeTokenKey(tokenIdentifier, tokenNonce)
			tokenBalance := instance.Balance.Value
			var uris [][]byte
			for _, jsonUri := range instance.Uris.Values {
				uris = append(uris, jsonUri.Value)
			}
			tokenData := &esdt.ESDigitalToken{
				Value:      tokenBalance,
				Type:       uint32(core.Fungible),
				Properties: MakeESDTUserMetadataBytes(isFrozen),
				TokenMetaData: &esdt.MetaData{
					Name:       []byte{},
					Nonce:      tokenNonce,
					Creator:    instance.Creator.Value,
					Royalties:  uint32(instance.Royalties.Value),
					Hash:       instance.Hash.Value,
					URIs:       uris,
					Attributes: instance.Attributes.Value,
				},
			}
			err := setTokenDataByKey(tokenKey, tokenData, destination)
			if err != nil {
				return err
			}
		}
		err := SetLastNonce(tokenIdentifier, scenESDTData.LastNonce.Value, destination)
		if err != nil {
			return err
		}
		err = SetTokenRolesAsStrings(tokenIdentifier, scenESDTData.Roles, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetTokenData sets the ESDT information related to a token into the storage of the account.
func setTokenDataByKey(tokenKey []byte, tokenData *esdt.ESDigitalToken, destination map[string][]byte) error {
	marshaledData, err := esdtDataMarshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}
	destination[string(tokenKey)] = marshaledData
	return nil
}

// SetTokenData sets the token data
func SetTokenData(tokenIdentifier []byte, nonce uint64, tokenData *esdt.ESDigitalToken, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return setTokenDataByKey(tokenKey, tokenData, destination)
}

// SetTokenRoles sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRoles(tokenIdentifier []byte, roles [][]byte, destination map[string][]byte) error {
	tokenRolesKey := makeTokenRolesKey(tokenIdentifier)
	tokenRolesData := &esdt.ESDTRoles{
		Roles: roles,
	}

	marshaledData, err := esdtDataMarshalizer.Marshal(tokenRolesData)
	if err != nil {
		return err
	}

	destination[string(tokenRolesKey)] = marshaledData
	return nil
}

// SetTokenRolesAsStrings sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRolesAsStrings(tokenIdentifier []byte, rolesAsStrings []string, destination map[string][]byte) error {
	roles := make([][]byte, len(rolesAsStrings))
	for i := 0; i < len(roles); i++ {
		roles[i] = []byte(rolesAsStrings[i])
	}

	return SetTokenRoles(tokenIdentifier, roles, destination)
}

// SetLastNonce writes the last nonce of a specified ESDT into the storage.
func SetLastNonce(tokenIdentifier []byte, lastNonce uint64, destination map[string][]byte) error {
	tokenNonceKey := makeLastNonceKey(tokenIdentifier)
	nonceBytes := big.NewInt(0).SetUint64(lastNonce).Bytes()
	destination[string(tokenNonceKey)] = nonceBytes
	return nil
}

// SetTokenBalance sets the ESDT balance of the account, specified by the token
// key.
func SetTokenBalance(tokenIdentifier []byte, nonce uint64, balance *big.Int, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	tokenData, err := getTokenDataByKey(tokenKey, destination, make(map[string][]byte))
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return errNegativeValue
	}

	tokenData.Value = balance
	return setTokenDataByKey(tokenKey, tokenData, destination)
}
