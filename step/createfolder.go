package step

import (
	"fmt"
	"github.com/ChristianHuff-DEV/reapy/model"
	log "github.com/sirupsen/logrus"
	"os"
)

// KindCreateFolder defines the name for a create folder step in the config file
const KindCreateFolder = "CreateFolder"

// CreateFolder is a step that is used to create a folder
type CreateFolder struct {
	model.RunnableStep
	Path string
}

// FromConfig create the struct representation of a step creating a folder
func (createFolder *CreateFolder) FromConfig(configYaml map[string]interface{}) error {
	createFolder.Kind = KindCreateFolder
	if description, ok := configYaml["Description"]; ok {
		createFolder.Description = description.(string)
	}
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	createFolder.Path = preferencesYaml["Path"].(string)
	return nil
}

// GetKind returns the type this step is of
func (createFolder CreateFolder) GetKind() string {
	return createFolder.Kind
}

// GetDescription returns a summary of what this step does
func (createFolder CreateFolder) GetDescription() string {
	return createFolder.Description
}

// Execute create a folder at a defined path
func (createFolder CreateFolder) Execute() model.Result {
	fmt.Println(createFolder.Description)
	log.Printf("Creating %s", createFolder.Path)
	return create(createFolder.Path)
}

// Create makes a new folder at the given path
func create(path string) (result model.Result) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = "folder created"
	return result
}
