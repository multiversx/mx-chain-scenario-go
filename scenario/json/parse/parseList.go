package scenjsonparse

import (
	"errors"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	mj "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

func (p *Parser) processStringList(obj interface{}) ([]string, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []string
	for _, elemRaw := range listRaw.AsList() {
		strVal, err := p.parseString(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, strVal)
	}
	return result, nil
}

func (p *Parser) parseValueList(obj interface{}) (mj.JSONValueList, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return mj.JSONValueList{}, errors.New("not a JSON list")
	}
	var result []mj.JSONBytesFromString
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processStringAsByteArray(elemRaw)
		if err != nil {
			return mj.JSONValueList{}, err
		}
		result = append(result, ba)
	}
	return mj.JSONValueList{
		Values: result,
	}, nil
}

func (p *Parser) parseSubTreeList(obj interface{}) ([]mj.JSONBytesFromTree, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []mj.JSONBytesFromTree
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processSubTreeAsByteArray(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, ba)
	}
	return result, nil
}

func (p *Parser) parseCheckValueList(obj oj.OJsonObject) (mj.JSONCheckValueList, error) {
	if IsStar(obj) {
		return mj.JSONCheckValueListStar(), nil
	}

	listRaw, listOk := obj.(*oj.OJsonList)
	if listOk {
		return p.parseCheckValueJSONList(listRaw)
	}

	if !p.AllowSingleValueInCheckValueList {
		return mj.JSONCheckValueList{}, errors.New("not a JSON list")
	}

	singleValue, err := p.parseCheckBytes(obj)
	if err != nil {
		return mj.JSONCheckValueList{}, err
	}

	if singleValue.OriginalEmpty() {
		// "" becomes [] instead of [""]
		return mj.JSONCheckValueList{
			Values: []mj.JSONCheckBytes{},
		}, nil
	}

	return mj.JSONCheckValueList{
		Values: []mj.JSONCheckBytes{singleValue},
	}, nil
}

func (p *Parser) parseCheckValueJSONList(listRaw *oj.OJsonList) (mj.JSONCheckValueList, error) {
	var values []mj.JSONCheckBytes
	for _, elemRaw := range listRaw.AsList() {
		checkBytes, err := p.parseCheckBytes(elemRaw)
		if err != nil {
			return mj.JSONCheckValueList{}, err
		}
		values = append(values, checkBytes)
	}
	return mj.JSONCheckValueList{
		Values: values,
	}, nil
}
