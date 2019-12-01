package step

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
)

type Download struct {
	model.RunnableStep
	DownloadURL  string
	DownloadPath string
}

func (this Download) GetKind() string {
	return this.Kind
}

func (this Download) GetDescription() string {
	return this.Description
}

func (this Download) Execute() model.Result {

	// Extract the filename from the last part of the URL (everything after the last "/")
	fileName := this.DownloadURL[strings.LastIndex(this.DownloadURL, "/")+1:]

	log.Print(this.DownloadURL)
	log.Print(fileName)

	if err := DownloadFile(this.DownloadPath+"/"+fileName, this.DownloadURL); err != nil {
		panic(err)
	}

	return model.Result{WasSuccessful: true}
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