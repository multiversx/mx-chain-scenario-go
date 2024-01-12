package util

import (
	"math/big"
	"testing"

	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	"github.com/stretchr/testify/require"
)

func Test_CreateMultiTransferData_SingleTransfer(t *testing.T) {
	esdtTransfers := make([]*scenmodel.ESDTTxData, 0)
	esdtTransfer := &scenmodel.ESDTTxData{
		Nonce:           scenmodel.JSONUint64{Value: 0},
		TokenIdentifier: scenmodel.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           scenmodel.JSONBigInt{Value: big.NewInt(100)},
	}
	esdtTransfers = append(esdtTransfers, esdtTransfer)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		esdtTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiESDTNFTTransfer@726563656976657241646472657373@01@544f4b312d616263646566@@64@66756e6374696f6e31@61726731@61726732", string(data))
}

func Test_CreateMultiTransferData_MultipleTransfers(t *testing.T) {
	esdtTransfers := make([]*scenmodel.ESDTTxData, 0)
	esdtTransfer1 := &scenmodel.ESDTTxData{
		Nonce:           scenmodel.JSONUint64{Value: 2},
		TokenIdentifier: scenmodel.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           scenmodel.JSONBigInt{Value: big.NewInt(100)},
	}
	esdtTransfer2 := &scenmodel.ESDTTxData{
		Nonce:           scenmodel.JSONUint64{Value: 14},
		TokenIdentifier: scenmodel.JSONBytesFromString{Value: []byte("TOK2-abcdef")},
		Value:           scenmodel.JSONBigInt{Value: big.NewInt(396)},
	}

	esdtTransfers = append(esdtTransfers, esdtTransfer1, esdtTransfer2)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		esdtTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiESDTNFTTransfer@726563656976657241646472657373@02@544f4b312d616263646566@02@64@544f4b322d616263646566@0e@018c@66756e6374696f6e31@61726731@61726732", string(data))
}
