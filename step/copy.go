package step

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindCopy is the name used to define this kind of step in the config file
const KindCopy = "Copy"

// Copy represents the struct of a step type that can be used to copy a file or folder
type Copy struct {
	model.RunnableStep
	Source      string
	Destination string
}

// GetKind returns the type of this step
func (copy Copy) GetKind() string {
	return copy.Kind
}

// GetDescription gives a short summary of what this step does
func (copy Copy) GetDescription() string {
	return copy.Description
}

// FromConfig creates the copy struct based on the provided config data
func (copy *Copy) FromConfig(stepConfig map[string]interface{}) error {
	copy.Kind = KindCopy

	if description, ok := stepConfig["Description"].(string); ok {
		copy.Description = description
	}

	if preferencesYaml, ok := stepConfig["Preferences"].(map[string]interface{}); ok {
		// Read source preference
		if source, ok := preferencesYaml["Source"].(string); ok {
			copy.Source = source
		} else {
			return fmt.Errorf("preference \"Source\" (string) must be set for %s step", copy.GetKind())
		}
		// Read destination preference
		if destination, ok := preferencesYaml["Destination"].(string); ok {
			copy.Destination = destination
		} else {
			return fmt.Errorf("preference \"Destination\" (string) must be set for %s step", copy.GetKind())
		}
	} else {
		return fmt.Errorf("preferences must be set for %s step", copy.GetKind())
	}
	return nil
}

// Execute determins if a file or folder has to be copied an triggers the apropiate copy operation
func (copy Copy) Execute() (result model.Result) {
	m := fmt.Sprintf("Copy %s to %s", copy.Source, copy.Destination)
	fmt.Println(m)
	log.Print(m)

	var err error

	fileInfo, err := os.Stat(copy.Source)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		err = dir(copy.Source, copy.Destination)
	case mode.IsRegular():
		err = file(copy.Source, copy.Destination)
	}

	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = fmt.Sprintf("Copied %s to %s", copy.Source, copy.Destination)
	return result
}

func dir(src string, dst string) error {
	var err error
	var filesInSourceFolder []os.FileInfo
	var sourceFolder os.FileInfo

	if sourceFolder, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, sourceFolder.Mode()); err != nil {
		return err
	}

	if filesInSourceFolder, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fileInSourceFolder := range filesInSourceFolder {
		sourcePath := path.Join(src, fileInSourceFolder.Name())
		destinationPath := path.Join(dst, fileInSourceFolder.Name())

		if fileInSourceFolder.IsDir() {
			if err = dir(sourcePath, destinationPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = file(sourcePath, destinationPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func file(src, dst string) error {
	var err error
	var sourceFile *os.File
	var destinationFile *os.File
	var srcinfo os.FileInfo

	if sourceFile, err = os.Open(src); err != nil {
		return err
	}
	defer sourceFile.Close()

	if destinationFile, err = os.Create(dst); err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
