package scenjsonparse

import (
	"errors"
	"fmt"

	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) (scenmodel.LogList, error) {
	if IsStar(logsRaw) {
		return scenmodel.LogList{
			IsUnspecified: false,
			IsStar:        true,
		}, nil
	}

	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return scenmodel.LogList{}, errors.New("unmarshalled logs list is not a list")
	}
	result := scenmodel.LogList{
		IsUnspecified:    false,
		IsStar:           false,
		MoreAllowedAtEnd: false,
		List:             nil,
	}
	var err error
	for _, logRaw := range logList.AsList() {
		switch logItem := logRaw.(type) {
		case *oj.OJsonString:
			if logItem.Value == "+" {
				result.MoreAllowedAtEnd = true
			} else {
				return scenmodel.LogList{}, errors.New("unmarshalled log entry is an invalid string")
			}
		case *oj.OJsonMap:
			if result.MoreAllowedAtEnd {
				return scenmodel.LogList{}, errors.New("log entry ")
			}

			logEntry := scenmodel.LogEntry{}
			for _, kvp := range logItem.OrderedKV {
				switch kvp.Key {
				case "address":
					logEntry.Address, err = p.parseCheckBytes(kvp.Value)
					if err != nil {
						return scenmodel.LogList{}, fmt.Errorf("invalid log address: %w", err)
					}
				case "endpoint":
					logEntry.Endpoint, err = p.parseCheckBytes(kvp.Value)
					if err != nil {
						return scenmodel.LogList{}, fmt.Errorf("invalid log identifier: %w", err)
					}
				case "topics":
					logEntry.Topics, err = p.parseCheckValueList(kvp.Value)
					if err != nil {
						return scenmodel.LogList{}, fmt.Errorf("invalid log entry topics: %w", err)
					}
				case "data":
					logEntry.Data, err = p.parseCheckValueList(kvp.Value)
					if err != nil {
						return scenmodel.LogList{}, fmt.Errorf("invalid log data: %w", err)
					}
				default:
					return scenmodel.LogList{}, fmt.Errorf("unknown log field: %s", kvp.Key)
				}
			}
			result.List = append(result.List, &logEntry)
		default:
			return scenmodel.LogList{}, errors.New("log entry should be either string or object")
		}
	}

	return result, nil
}
