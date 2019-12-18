package step

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindCopyFile defines the name for a copy step in the config file
const KindCopyFile = "CopyFile"

// CopyFile takes a file and copies it from one location to another
type CopyFile struct {
	model.RunnableStep
	// Source is the path to the file to copy it from
	Source string
	// Destination is the path to the file to copy it to
	Destination string
}

// GetKind returns the kind this step represents
func (copyFile CopyFile) GetKind() string {
	return copyFile.Kind
}

// GetDescription gives a description of what this step does
func (copyFile CopyFile) GetDescription() string {
	return copyFile.Description
}

// FromConfig creates the CopyFile struct from the given config
func (copyFile *CopyFile) FromConfig(stepConfig map[string]interface{}) error {
	copyFile.Kind = KindCopyFile

	if description, ok := stepConfig["Description"].(string); ok {
		copyFile.Description = description
	}

	if preferencesYaml, ok := stepConfig["Preferences"].(map[string]interface{}); ok {
		// Read source preference
		if source, ok := preferencesYaml["Source"].(string); ok {
			copyFile.Source = source
		} else {
			return fmt.Errorf("preference \"Source\" (string) must be set for %s step", copyFile.GetKind())
		}
		// Read destination preference
		if destination, ok := preferencesYaml["Destination"].(string); ok {
			copyFile.Destination = destination
		} else {
			return fmt.Errorf("preference \"Destination\" (string) must be set for %s step", copyFile.GetKind())
		}
	} else {
		return fmt.Errorf("preferences must be set for %s step", copyFile.GetKind())
	}
	return nil
}

// Execute copies a file/folder from one location to another
func (copyFile CopyFile) Execute() (result model.Result) {
	m := fmt.Sprintf("CopyFile %s to %s", copyFile.Source, copyFile.Destination)
	fmt.Println(m)
	log.Print(m)

	source, err := os.Open(copyFile.Source)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}
	defer source.Close()

	destination, err := os.Create(copyFile.Destination)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = fmt.Sprintf("Copied %s to %s", copyFile.Source, copyFile.Destination)
	return result
}
