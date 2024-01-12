package scenjsonparse

import (
	"errors"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
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

func (p *Parser) parseValueList(obj interface{}) (scenmodel.JSONValueList, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return scenmodel.JSONValueList{}, errors.New("not a JSON list")
	}
	var result []scenmodel.JSONBytesFromString
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processStringAsByteArray(elemRaw)
		if err != nil {
			return scenmodel.JSONValueList{}, err
		}
		result = append(result, ba)
	}
	return scenmodel.JSONValueList{
		Values: result,
	}, nil
}

func (p *Parser) parseSubTreeList(obj interface{}) ([]scenmodel.JSONBytesFromTree, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []scenmodel.JSONBytesFromTree
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processSubTreeAsByteArray(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, ba)
	}
	return result, nil
}

func (p *Parser) parseCheckValueList(obj oj.OJsonObject) (scenmodel.JSONCheckValueList, error) {
	if IsStar(obj) {
		return scenmodel.JSONCheckValueListStar(), nil
	}

	listRaw, listOk := obj.(*oj.OJsonList)
	if listOk {
		return p.parseCheckValueJSONList(listRaw)
	}

	if !p.AllowSingleValueInCheckValueList {
		return scenmodel.JSONCheckValueList{}, errors.New("not a JSON list")
	}

	singleValue, err := p.parseCheckBytes(obj)
	if err != nil {
		return scenmodel.JSONCheckValueList{}, err
	}

	if singleValue.OriginalEmpty() {
		// "" becomes [] instead of [""]
		return scenmodel.JSONCheckValueList{
			Values: []scenmodel.JSONCheckBytes{},
		}, nil
	}

	return scenmodel.JSONCheckValueList{
		Values: []scenmodel.JSONCheckBytes{singleValue},
	}, nil
}

func (p *Parser) parseCheckValueJSONList(listRaw *oj.OJsonList) (scenmodel.JSONCheckValueList, error) {
	var values []scenmodel.JSONCheckBytes
	for _, elemRaw := range listRaw.AsList() {
		checkBytes, err := p.parseCheckBytes(elemRaw)
		if err != nil {
			return scenmodel.JSONCheckValueList{}, err
		}
		values = append(values, checkBytes)
	}
	return scenmodel.JSONCheckValueList{
		Values: values,
	}, nil
}
