package step

import "github.com/ChristianHuff-DEV/reapy/model"

import "os"

import "log"

// KindCreateFolder defines the name for a create folder step in the config file
const KindCreateFolder = "CreateFolder"

// CreateFolder is a step that is used to create a folder
type CreateFolder struct {
	model.RunnableStep
	Path string
}

// FromConfig create the struct representation of a step creating a folder
func (createFolder *CreateFolder) FromConfig(configYaml map[string]interface{}) {
	createFolder.Kind = KindCreateFolder
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	createFolder.Path = preferencesYaml["Path"].(string)
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
	return create(createFolder.Path)
}

// Create makes a new folder at the given path
func create(path string) model.Result {
	log.Printf("Creating %s", path)
	//Attempt to create the directory and ignore any issues
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		log.Print(err)
	}
	return model.Result{WasSuccessful: true}
}
