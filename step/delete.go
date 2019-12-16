package step

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindDelete defines the name for a delete step in the config file
const KindDelete = "Delete"

// Delete is a step used to delete a file/folder
type Delete struct {
	model.RunnableStep
	Path string
}

// GetKind returns the type this step is of
func (delete Delete) GetKind() string {
	return delete.Kind
}

// GetDescription returns a summary of what this step does
func (delete Delete) GetDescription() string {
	return delete.Description
}

// FromConfig creates the struct representation of a step deleting a file/folder
func (delete *Delete) FromConfig(configYaml map[string]interface{}) error {
	delete.Kind = KindDelete
	if description, ok := configYaml["Description"]; ok {
		delete.Description = description.(string)
	}
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	delete.Path = preferencesYaml["Path"].(string)
	return nil
}

// Execute trigger the deletion of a file/folder
func (delete Delete) Execute() model.Result {
	fmt.Println(delete.Description)
	log.Printf("Deleting %s", delete.Path)
	return deletePath(delete.Path)
}

// delete removes the file/folder at the given path
func deletePath(path string) (result model.Result) {
	f, err := os.Stat(path)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	if f.IsDir() {
		err := os.RemoveAll(path)
		if err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
	} else {
		err := os.Remove(path)
		if err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
	}

	result.WasSuccessful = true
	result.Message = path + " deleted"
	return result
}
