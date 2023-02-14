package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/multiversx/mx-chain-core-go/core"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	txDataBuilder "github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
)

// CreateMultiTransferData builds data for a multiTransferESDT
func CreateMultiTransferData(to []byte, esdtData []*mj.ESDTTxData, endpointName string, arguments [][]byte) []byte {

	multiTransferData := make([]byte, 0)
	multiTransferData = append(multiTransferData, []byte(core.BuiltInFunctionMultiESDTNFTTransfer)...)
	tdb := txDataBuilder.NewBuilder()
	tdb.Bytes(to)
	tdb.Int(len(esdtData))

	for _, esdtDataTransfer := range esdtData {
		tdb.Bytes(esdtDataTransfer.TokenIdentifier.Value)
		tdb.Int64(int64(esdtDataTransfer.Nonce.Value))
		tdb.BigInt(esdtDataTransfer.Value.Value)
	}

	if len(endpointName) > 0 {
		tdb.Str(endpointName)

		for _, arg := range arguments {
			tdb.Bytes(arg)
		}
	}
	multiTransferData = append(multiTransferData, tdb.ToBytes()...)
	return multiTransferData
}

// GetSCCode retrieves the bytecode of a WASM module from a file
func GetSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}

	return code
}
