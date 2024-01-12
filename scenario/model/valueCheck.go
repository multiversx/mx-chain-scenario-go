package scenmodel

import (
	"bytes"
	"math/big"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
)

// JSONCheckBytes holds a byte slice condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBytes struct {
	Value       []byte
	IsStar      bool
	Original    oj.OJsonObject
	Unspecified bool
}

// JSONCheckBytesUnspecified yields JSONCheckBytes that check that value is empty.
func JSONCheckBytesUnspecified() JSONCheckBytes {
	return JSONCheckBytes{
		Value:       []byte{},
		IsStar:      false,
		Original:    &oj.OJsonString{Value: ""},
		Unspecified: true,
	}
}

// JSONCheckBytesStar yields JSONCheckBytes explicit "*" value.
func JSONCheckBytesStar() JSONCheckBytes {
	return JSONCheckBytes{
		Value:       []byte{},
		IsStar:      true,
		Original:    &oj.OJsonString{Value: "*"},
		Unspecified: false,
	}
}

// JSONCheckBytesReconstructed creates a JSONCheckBytes without an original JSON source.
func JSONCheckBytesReconstructed(value []byte, originalString string) JSONCheckBytes {
	return JSONCheckBytes{
		Value:       value,
		IsStar:      false,
		Original:    &oj.OJsonString{Value: originalString},
		Unspecified: false,
	}
}

// OriginalEmpty returns true if original = "".
func (jcbytes JSONCheckBytes) OriginalEmpty() bool {
	if str, isStr := jcbytes.Original.(*oj.OJsonString); isStr {
		return len(str.Value) == 0
	}
	return false
}

// IsUnspecified yields true if the field was originally unspecified.
func (jcbytes JSONCheckBytes) IsUnspecified() bool {
	return jcbytes.Unspecified
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbytes JSONCheckBytes) Check(other []byte) bool {
	if jcbytes.IsStar {
		return true
	}
	return bytes.Equal(jcbytes.Value, other)
}

// JSONCheckBigInt holds a big int condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBigInt struct {
	Value       *big.Int
	IsStar      bool
	Original    string
	Unspecified bool
}

// JSONCheckBigIntUnspecified yields JSONCheckBigInt default "*" value.
func JSONCheckBigIntUnspecified() JSONCheckBigInt {
	return JSONCheckBigInt{
		Value:       big.NewInt(0),
		IsStar:      false,
		Original:    "",
		Unspecified: true,
	}
}

// IsUnspecified yields true if the field was originally unspecified.
func (jcbi JSONCheckBigInt) IsUnspecified() bool {
	return jcbi.Unspecified
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbi JSONCheckBigInt) Check(other *big.Int) bool {
	if jcbi.IsStar {
		return true
	}
	return jcbi.Value.Cmp(other) == 0
}

// JSONCheckUint64 holds a uint64 condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckUint64 struct {
	Value       uint64
	IsStar      bool
	Original    string
	Unspecified bool
}

// JSONCheckUint64Unspecified yields JSONCheckBigInt default "*" value.
func JSONCheckUint64Unspecified() JSONCheckUint64 {
	return JSONCheckUint64{
		Value:       0,
		IsStar:      false,
		Original:    "",
		Unspecified: true,
	}
}

// IsUnspecified yields true if the field was originally unspecified.
func (jcu JSONCheckUint64) IsUnspecified() bool {
	return jcu.Unspecified
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcu JSONCheckUint64) Check(other uint64) bool {
	if jcu.IsStar {
		return true
	}
	return jcu.Value == other
}

// CheckBool interprets own value as bool (true = anything > 0, false = 0),
// We are using JSONCheckUint64 for bool too so we don't create another type.
func (jcu JSONCheckUint64) CheckBool(other bool) bool {
	if jcu.IsStar {
		return true
	}
	return jcu.Value > 0 == other
}

// JSONCheckValueList represents a list of value checks, as expressed in JSON.
// TODO: add star for all values
type JSONCheckValueList struct {
	Values      []JSONCheckBytes
	IsStar      bool
	Unspecified bool
}

// JSONCheckValueListUnspecified yields JSONCheckBytesList empty value.
func JSONCheckValueListUnspecified() JSONCheckValueList {
	return JSONCheckValueList{
		Values:      nil,
		IsStar:      false,
		Unspecified: true,
	}
}

// JSONCheckValueListStar yields the "*" value
func JSONCheckValueListStar() JSONCheckValueList {
	return JSONCheckValueList{
		Values:      nil,
		IsStar:      true,
		Unspecified: false,
	}
}

// IsUnspecified yields true if the field was originally unspecified.
func (jcbl JSONCheckValueList) IsUnspecified() bool {
	return jcbl.Unspecified
}

// CheckList compares expected value with a list of values.
func (jcbl JSONCheckValueList) CheckList(other [][]byte) bool {
	if jcbl.IsStar {
		return true
	}
	if len(jcbl.Values) != len(other) {
		return false
	}
	for i, expected := range jcbl.Values {
		if !expected.Check(other[i]) {
			return false
		}
	}
	return true
}
