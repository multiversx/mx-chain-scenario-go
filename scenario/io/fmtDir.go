package scenio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var suffixes = []string{".scen.json", ".step.json", ".steps.json"}

func shouldFormatFile(path string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}

func FormatAllInFolder(path string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if shouldFormatFile(filePath) {
			fmt.Printf("Formatting: %s\n", filePath)
			upgradeScenariosFile(filePath)
		}
		return nil
	})
	return err
}

func upgradeScenariosFile(filePath string) {
	scenario, err := ParseScenariosScenarioDefaultParser(filePath)
	if err == nil {
		_ = WriteScenariosScenario(scenario, filePath)
	} else {
		fmt.Printf("Error upgrading: %s\n", err.Error())
	}
}
