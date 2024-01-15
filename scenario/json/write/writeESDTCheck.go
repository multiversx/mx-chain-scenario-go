package scenjsonwrite

import (
	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

func checkESDTDataToOJ(esdtItems []*scenmodel.CheckESDTData, moreESDTTokensAllowed bool) *oj.OJsonMap {
	esdtItemsOJ := oj.NewMap()
	for _, esdtItem := range esdtItems {
		esdtItemsOJ.Put(esdtItem.TokenIdentifier.Original, checkESDTItemToOJ(esdtItem))
	}
	if moreESDTTokensAllowed {
		esdtItemsOJ.Put("+", stringToOJ(""))
	}
	return esdtItemsOJ
}

func checkESDTItemToOJ(esdtItem *scenmodel.CheckESDTData) oj.OJsonObject {
	if isCompactCheckESDT(esdtItem) {
		return checkBigIntToOJ(esdtItem.Instances[0].Balance)
	}

	esdtItemOJ := oj.NewMap()

	// instances
	if len(esdtItem.Instances) > 0 {
		var convertedList []oj.OJsonObject
		for _, esdtInstance := range esdtItem.Instances {
			esdtInstanceOJ := oj.NewMap()
			appendCheckESDTInstanceToOJ(esdtInstance, esdtInstanceOJ)
			convertedList = append(convertedList, esdtInstanceOJ)
		}
		instancesOJList := oj.OJsonList(convertedList)
		esdtItemOJ.Put("instances", &instancesOJList)
	}

	if len(esdtItem.LastNonce.Original) > 0 {
		esdtItemOJ.Put("lastNonce", checkUint64ToOJ(esdtItem.LastNonce))
	}

	// roles
	if len(esdtItem.Roles) > 0 {
		var convertedList []oj.OJsonObject
		for _, roleStr := range esdtItem.Roles {
			convertedList = append(convertedList, &oj.OJsonString{Value: roleStr})
		}
		rolesOJList := oj.OJsonList(convertedList)
		esdtItemOJ.Put("roles", &rolesOJList)
	}
	if len(esdtItem.Frozen.Original) > 0 {
		esdtItemOJ.Put("frozen", checkUint64ToOJ(esdtItem.Frozen))
	}

	return esdtItemOJ
}

func appendCheckESDTInstanceToOJ(esdtInstance *scenmodel.CheckESDTInstance, targetOj *oj.OJsonMap) {
	targetOj.Put("nonce", uint64ToOJ(esdtInstance.Nonce))

	if len(esdtInstance.Balance.Original) > 0 {
		targetOj.Put("balance", checkBigIntToOJ(esdtInstance.Balance))
	}
	if !esdtInstance.Creator.Unspecified && len(esdtInstance.Creator.Value) > 0 {
		targetOj.Put("creator", checkBytesToOJ(esdtInstance.Creator))
	}
	if !esdtInstance.Royalties.Unspecified && len(esdtInstance.Royalties.Original) > 0 {
		targetOj.Put("royalties", checkUint64ToOJ(esdtInstance.Royalties))
	}
	if !esdtInstance.Hash.Unspecified && len(esdtInstance.Hash.Value) > 0 {
		targetOj.Put("hash", checkBytesToOJ(esdtInstance.Hash))
	}
	if !esdtInstance.Uris.IsUnspecified() {
		targetOj.Put("uri", checkValueListToOJ(esdtInstance.Uris))
	}
	if !esdtInstance.Attributes.Unspecified && len(esdtInstance.Attributes.Value) > 0 {
		targetOj.Put("attributes", checkBytesToOJ(esdtInstance.Attributes))
	}
}

func isCompactCheckESDT(esdtItem *scenmodel.CheckESDTData) bool {
	if len(esdtItem.Instances) != 1 {
		return false
	}
	if len(esdtItem.Instances[0].Nonce.Original) > 0 {
		return false
	}
	if len(esdtItem.Roles) > 0 {
		return false
	}
	if len(esdtItem.Frozen.Original) > 0 {
		return false
	}
	return true
}
