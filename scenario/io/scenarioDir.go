package scenio

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/TwiN/go-color"
)

// RunAllJSONScenariosInDirectory walks directory, parses and prepares all json scenarios,
// then calls ScenarioRunner for each of them.
func (r *ScenarioController) RunAllJSONScenariosInDirectory(
	generalTestPath string,
	specificTestPath string,
	allowedSuffix string,
	excludedFilePatterns []string,
	options *RunScenarioOptions) error {

	mainDirPath := path.Join(generalTestPath, specificTestPath)
	var nrPassed, nrFailed, nrSkipped int

	err := filepath.Walk(mainDirPath, func(testFilePath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(testFilePath, allowedSuffix) {
			fmt.Printf("Scenario: %s ... ", shortenTestPath(testFilePath, generalTestPath))
			if isExcluded(excludedFilePatterns, testFilePath, generalTestPath) {
				nrSkipped++
				fmt.Printf("  %s\n", color.Ize(color.Yellow, "skip"))
			} else {
				r.Executor.Reset()
				r.RunsNewTest = true
				testErr := r.RunSingleJSONScenario(testFilePath, options)
				if testErr == nil {
					nrPassed++
					fmt.Printf("  %s\n", color.Ize(color.Green, "ok"))
				} else {
					nrFailed++
					fmt.Printf("  %s %s\n", color.Ize(color.Red, "FAIL:"), testErr.Error())
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Done. Passed: %d. Failed: %d. Skipped: %d.\n", nrPassed, nrFailed, nrSkipped)
	if nrFailed > 0 {
		return errors.New("some tests failed")
	}

	return nil
}

func isExcluded(excludedFilePatterns []string, testPath string, generalTestPath string) bool {
	for _, et := range excludedFilePatterns {
		excludedFullPath := path.Join(generalTestPath, et)
		match, err := filepath.Match(excludedFullPath, testPath)
		if err != nil {
			panic(err)
		}
		if match {
			return true
		}
	}
	return false
}

func shortenTestPath(path string, generalTestPath string) string {
	if strings.HasPrefix(path, generalTestPath+"/") {
		return path[len(generalTestPath)+1:]
	}
	return path
}
