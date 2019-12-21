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

// KindCopyFolder defines the name for a copy step in the config file
const KindCopyFolder = "CopyFolder"

// CopyFolder copies the all the content of one folder to another folder.
//
// If the destination file already exists it will be merged into the existing folder. All existing files will be overwritten with the new files.
type CopyFolder struct {
	model.RunnableStep
	// Source is the path to the folder to be copied
	Source string
	// Destination is the path to copy the folder to
	Destination string
}

// GetKind returns the kind this step represents
func (copyFolder CopyFolder) GetKind() string {
	return copyFolder.Kind
}

// GetDescription gives a description of what this step does
func (copyFolder CopyFolder) GetDescription() string {
	return copyFolder.Description
}

// FromConfig creates the CopyFolder struct from the given config
func (copyFolder *CopyFolder) FromConfig(stepConfig map[string]interface{}) error {
	copyFolder.Kind = KindCopyFolder

	if description, ok := stepConfig["Description"].(string); ok {
		copyFolder.Description = description
	}

	if preferencesYaml, ok := stepConfig["Preferences"].(map[string]interface{}); ok {
		// Read source preference
		if source, ok := preferencesYaml["Source"].(string); ok {
			copyFolder.Source = source
		} else {
			return fmt.Errorf("preference \"Source\" (string) must be set for %s step", copyFolder.GetKind())
		}
		// Read destination preference
		if destination, ok := preferencesYaml["Destination"].(string); ok {
			copyFolder.Destination = destination
		} else {
			return fmt.Errorf("preference \"Destination\" (string) must be set for %s step", copyFolder.GetKind())
		}
	} else {
		return fmt.Errorf("preferences must be set for %s step", copyFolder.GetKind())
	}
	return nil
}

// Execute copies a folder from one location to another
func (copyFolder CopyFolder) Execute() (result model.Result) {
	m := fmt.Sprintf("Copy folder from %s to %s", copyFolder.Source, copyFolder.Destination)
	fmt.Println(m)
	log.Print(m)

	err := Dir(copyFolder.Source, copyFolder.Destination)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = fmt.Sprintf("Copied %s to %s", copyFolder.Source, copyFolder.Destination)
	return result
}

// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
