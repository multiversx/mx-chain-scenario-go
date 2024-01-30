package scenjsonparse

import (
	"errors"
	"math/big"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

type bigIntParseFormat int

const (
	bigIntSignedBytes bigIntParseFormat = iota
	bigIntUnsignedBytes
)

func (p *Parser) processCheckBigInt(obj oj.OJsonObject, format bigIntParseFormat) (scenmodel.JSONCheckBigInt, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return scenmodel.JSONCheckBigInt{
			Value:    nil,
			IsStar:   true,
			Original: "*"}, nil
	}

	jbi, err := p.processBigInt(obj, format)
	if err != nil {
		return scenmodel.JSONCheckBigInt{}, err
	}
	return scenmodel.JSONCheckBigInt{
		Value:    jbi.Value,
		IsStar:   false,
		Original: jbi.Original,
	}, nil
}

func (p *Parser) processBigInt(obj oj.OJsonObject, format bigIntParseFormat) (scenmodel.JSONBigInt, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return scenmodel.JSONBigInt{}, err
	}

	bi, err := p.parseBigInt(strVal, format)
	return scenmodel.JSONBigInt{
		Value:    bi,
		Original: strVal,
	}, err
}

func (p *Parser) parseBigInt(strRaw string, format bigIntParseFormat) (*big.Int, error) {
	bytes, err := p.ExprInterpreter.InterpretString(strRaw)
	if err != nil {
		return nil, err
	}
	switch format {
	case bigIntSignedBytes:
		return twos.FromBytes(bytes), nil
	case bigIntUnsignedBytes:
		return big.NewInt(0).SetBytes(bytes), nil
	default:
		return nil, errors.New("unknown format requested")
	}
}

func (p *Parser) processCheckUint64(obj oj.OJsonObject) (scenmodel.JSONCheckUint64, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return scenmodel.JSONCheckUint64{
			Value:    0,
			IsStar:   true,
			Original: "*"}, nil
	}

	ju, err := p.processUint64(obj)
	if err != nil {
		return scenmodel.JSONCheckUint64{}, err
	}
	return scenmodel.JSONCheckUint64{
		Value:    ju.Value,
		IsStar:   false,
		Original: ju.Original}, nil

}

func (p *Parser) processUint64(obj oj.OJsonObject) (scenmodel.JSONUint64, error) {
	bi, err := p.processBigInt(obj, bigIntUnsignedBytes)
	if err != nil {
		return scenmodel.JSONUint64{}, err
	}

	if bi.Value == nil || !bi.Value.IsUint64() {
		return scenmodel.JSONUint64{}, errors.New("value is not uint64")
	}

	return scenmodel.JSONUint64{
		Value:    bi.Value.Uint64(),
		Original: bi.Original}, nil
}

func (p *Parser) parseCheckBytes(obj oj.OJsonObject) (scenmodel.JSONCheckBytes, error) {
	if IsStar(obj) {
		// "*" means any value, skip checking it
		return scenmodel.JSONCheckBytesStar(), nil
	}

	jb, err := p.processSubTreeAsByteArray(obj)
	if err != nil {
		return scenmodel.JSONCheckBytes{}, err
	}
	return scenmodel.JSONCheckBytes{
		Value:    jb.Value,
		IsStar:   false,
		Original: jb.Original,
	}, nil
}

func (p *Parser) processStringAsByteArray(obj oj.OJsonObject) (scenmodel.JSONBytesFromString, error) {
	strVal, err := p.parseString(obj)
	if err != nil {
		return scenmodel.JSONBytesFromString{}, err
	}
	result, err := p.ExprInterpreter.InterpretString(strVal)
	return scenmodel.NewJSONBytesFromString(result, strVal), err
}

func (p *Parser) processSubTreeAsByteArray(obj oj.OJsonObject) (scenmodel.JSONBytesFromTree, error) {
	value, err := p.ExprInterpreter.InterpretSubTree(obj)
	return scenmodel.JSONBytesFromTree{
		Value:    value,
		Original: obj,
	}, err
}

func (p *Parser) parseString(obj oj.OJsonObject) (string, error) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", errors.New("not a string value")
	}
	return str.Value, nil
}

func (p *Parser) parseBool(obj oj.OJsonObject) (bool, error) {
	value, isBool := obj.(*oj.OJsonBool)
	if !isBool {
		return false, errors.New("not a bool value")
	}
	return bool(*value), nil
}

// IsStar returns whether check object is othe form "*".
func IsStar(obj oj.OJsonObject) bool {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return false
	}
	return str.Value == "*"
}
