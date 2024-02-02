package scenjsontest

import (
	"io"
	"os"
	"testing"

	fr "github.com/multiversx/mx-chain-scenario-go/scenario/expression/fileresolver"
	scenjparse "github.com/multiversx/mx-chain-scenario-go/scenario/json/parse"
	scenjwrite "github.com/multiversx/mx-chain-scenario-go/scenario/json/write"

	"github.com/stretchr/testify/require"
)

// only for this test
var vmType = []byte{'W', 'W'}

func loadExampleFile(path string) ([]byte, error) {
	// Open our jsonFile
	var jsonFile *os.File
	var err error
	jsonFile, err = os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		_ = jsonFile.Close()
	}()

	return io.ReadAll(jsonFile)
}

func TestWriteScenario(t *testing.T) {
	contents, err := loadExampleFile("example.scen.json")
	require.Nil(t, err)

	p := scenjparse.NewParser(
		fr.NewDefaultFileResolver().ReplacePath(
			"smart-contract.wasm",
			"exampleFile.txt"),
		vmType)

	scenario, parseErr := p.ParseScenarioFile(contents)
	require.Nil(t, parseErr)

	serialized := scenjwrite.ScenarioToJSONString(scenario)

	// good for debugging:
	_ = os.WriteFile("serialized.scen.json", []byte(serialized), 0644)

	require.Equal(t, contents, []byte(serialized))
}
