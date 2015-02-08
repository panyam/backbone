package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ResourceHandler struct {
}

func NewResourceHandler() (*ResourceHandler, error) {
	handler := ResourceHandler{}
	return &handler, nil
}

func (r *ResourceHandler) Download(url string, target string) error {
	// download
	log.Printf("Downloading %s -> %s\n", url, target)
	res, err := http.Get(url)
	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return fmt.Errorf("Could not download %s, status: %d", url, res.StatusCode)
	}

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	buffer := make([]byte, 1024)
	for n, err := 0, error(nil); n > 0 || err == nil; n, err = res.Body.Read(buffer) {
		if n == 0 {
			continue
		}
		_, writeErr := file.Write(buffer[0:n])
		if writeErr != nil {
			return err
		}
	}
	err = os.Chtimes(target, time.Now(), time.Now())
	if err != nil {
		log.Println("Error changing time stamp for ", target, err)
	}

	return nil
}

func (r *ResourceHandler) Upload(source string, url string) error {
	return nil
}
