package scenjsonparse

import (
	fr "github.com/multiversx/mx-chain-scenario-go/scenario/expression/fileresolver"
	ei "github.com/multiversx/mx-chain-scenario-go/scenario/expression/interpreter"
)

// Parser performs parsing of both json tests (older) and scenarios (new).
type Parser struct {
	ExprInterpreter                  ei.ExprInterpreter
	AllowEsdtTxLegacySyntax          bool
	AllowEsdtLegacySetSyntax         bool
	AllowEsdtLegacyCheckSyntax       bool
	AllowSingleValueInCheckValueList bool
}

// NewParser provides a new Parser instance.
func NewParser(fileResolver fr.FileResolver, vmType []byte) Parser {
	return Parser{
		ExprInterpreter: ei.ExprInterpreter{
			FileResolver: fileResolver,
			VMType:       vmType,
		},
		AllowEsdtTxLegacySyntax:          true,
		AllowEsdtLegacySetSyntax:         true,
		AllowEsdtLegacyCheckSyntax:       true,
		AllowSingleValueInCheckValueList: true,
	}
}
