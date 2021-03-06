package step

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
)

// KindDownload defines the name for a download step in the config file
const KindDownload = "Download"

// Download represents a step used to download a file from aj url
type Download struct {
	model.RunnableStep
	URL  string
	Path string
}

// GetKind returns the type this step is of
func (download Download) GetKind() string {
	return download.Kind
}

// GetDescription returns a text summarizing what this step does
func (download Download) GetDescription() string {
	return download.Description
}

// FromConfig creates the struct representation of a download step
func (download *Download) FromConfig(configYaml map[string]interface{}) error {
	download.Kind = KindDownload
	if description, ok := configYaml["Description"]; ok {
		download.Description = description.(string)
	}
	preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	download.URL = preferencesYaml["URL"].(string)
	download.Path = preferencesYaml["Path"].(string)
	return nil
}

// Execute downloads the file found at a given url
func (download Download) Execute() (result model.Result) {
	fmt.Println(download.Description)
	log.Printf("Downloading %s", download.URL)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	// Extract the filename from the last part of the URL (everything after the last "/")
	fileName := download.URL[strings.LastIndex(download.URL, "/")+1:]

	if err := DownloadFile(download.Path+"/"+fileName, download.URL); err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		s.Stop()
		return result
	}

	s.Stop()
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
