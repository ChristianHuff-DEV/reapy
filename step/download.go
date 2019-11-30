package step

import (
	"io"
	"net/http"
	"os"

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

func (this Download) Execute() {
	if err := DownloadFile("tomcat.zip", this.DownloadURL); err != nil {
		panic(err)
	}
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
