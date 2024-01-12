package scenio

import (
	"io"
	"os"
	"path/filepath"

	scenjparse "github.com/multiversx/mx-chain-scenario-go/scenario/json/parse"
	scenjwrite "github.com/multiversx/mx-chain-scenario-go/scenario/json/write"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// ParseScenariosScenario reads and parses a Scenarios scenario from a JSON file.
func ParseScenariosScenario(parser scenjparse.Parser, scenFilePath string) (*scenmodel.Scenario, error) {
	var err error
	scenFilePath, err = filepath.Abs(scenFilePath)
	if err != nil {
		return nil, err
	}

	// Open our jsonFile
	var jsonFile *os.File
	jsonFile, err = os.Open(scenFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		_ = jsonFile.Close()
	}()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	parser.ExprInterpreter.FileResolver.SetContext(scenFilePath)
	return parser.ParseScenarioFile(byteValue)
}

// ParseScenariosScenarioDefaultParser reads and parses a Scenarios scenario from a JSON file.
func ParseScenariosScenarioDefaultParser(scenFilePath string) (*scenmodel.Scenario, error) {
	parser := scenjparse.NewParser(NewDefaultFileResolver())
	parser.ExprInterpreter.FileResolver.SetContext(scenFilePath)
	return ParseScenariosScenario(parser, scenFilePath)
}

// WriteScenariosScenario exports a Scenarios scenario to a file, using the default formatting.
func WriteScenariosScenario(scenario *scenmodel.Scenario, toPath string) error {
	jsonString := scenjwrite.ScenarioToJSONString(scenario)

	err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(toPath, []byte(jsonString), 0644)
}
