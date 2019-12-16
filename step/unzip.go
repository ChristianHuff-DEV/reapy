package step

import (
	"archive/zip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
)

// KindUnzip defines the name for a unzip step in the config file
const KindUnzip = "Unzip"

// Unzip represents a task used to unzip a zip archive file
type Unzip struct {
	model.RunnableStep
	Source      string
	Destination string
}

// GetKind returns the type this step is of
func (unzip Unzip) GetKind() string {
	return unzip.Kind
}

// GetDescription returns a text summarizing what this step does
func (unzip Unzip) GetDescription() string {
	return unzip.Description
}

// FromConfig creates the struct representation of a unzip step
func (unzip *Unzip) FromConfig(configYaml map[string]interface{}) error {
	unzip.Kind = KindUnzip
	if description, ok := configYaml["Description"]; ok {
		unzip.Description = description.(string)
	}
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	unzip.Source = preferencesYaml["Source"].(string)
	unzip.Destination = preferencesYaml["Destination"].(string)
	return nil
}

// Execute unzips a file to a defined location
func (unzip Unzip) Execute() model.Result {
	fmt.Println(unzip.Description)
	log.Printf("Unzipping %s to %s", unzip.Source, unzip.Destination)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	result := unzipFile(unzip.Source, unzip.Destination)
	s.Stop()

	return result
}

// unzipFile extracts the given source file to the destination folder
func unzipFile(src, dest string) (result model.Result) {

	// Create a reader for the specified file
	reader, err := zip.OpenReader(src)
	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}
	defer reader.Close()

	// Range over the content of the archive
	for _, f := range reader.File {

		rc, err := f.Open()
		if err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		// Check if the current file points to a file or a folder
		if f.FileInfo().IsDir() {
			// Create the folder
			os.MkdirAll(path, os.ModePerm)
		} else {
			// Create the file
			var fileDir string
			if lastIndex := strings.LastIndex(path, string(os.PathSeparator)); lastIndex > -1 {
				fileDir = path[:lastIndex]
			}

			err = os.MkdirAll(fileDir, os.ModePerm)
			if err != nil {
				result.WasSuccessful = false
				result.Message = err.Error()
				return result
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				result.WasSuccessful = false
				result.Message = err.Error()
				return result
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				result.WasSuccessful = false
				result.Message = err.Error()
				return result
			}
		}
	}

	result.WasSuccessful = true
	result.Message = "unpacked"
	return result
}
