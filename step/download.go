package step

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindDownload defines the name for a download step in the config file
const KindDownload = "Download"

// Download represents a step used to download a file from aj url
type Download struct {
	model.RunnableStep
	URL  string
	Path string
}

// GetKind returns the type this task is of
func (download Download) GetKind() string {
	return download.Kind
}

// GetDescription returns a text summarizing what this step does
func (download Download) GetDescription() string {
	return download.Description
}

// FromConfig create the struct representation of a download step
func (download *Download) FromConfig(configYaml map[string]interface{}) {
	download.Kind = KindDownload
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	download.URL = preferencesYaml["URL"].(string)
	download.Path = preferencesYaml["Path"].(string)
}

// Execute downloads the file found at a given url
func (download Download) Execute() (result model.Result) {
	log.Printf("Downloading %s", download.URL)

	// Extract the filename from the last part of the URL (everything after the last "/")
	fileName := download.URL[strings.LastIndex(download.URL, "/")+1:]

	if err := DownloadFile(download.Path+"/"+fileName, download.URL); err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = "downloaded"
	return result
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
