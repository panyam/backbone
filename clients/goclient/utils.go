package goclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	// "strconv"
)

func MakeRequest(method, endpoint string, queryParams map[string]string, body io.Reader) (*http.Request, error) {
	if queryParams != nil {
		i := 0
		for key, value := range queryParams {
			if i == 0 {
				endpoint += "?"
			} else {
				endpoint += "&"
			}
			endpoint += key
			endpoint += "=" + value
			i++
		}
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Printf("%s %s Error: %s\n", method, endpoint, err)
		return nil, err
	}
	return req, nil
}

func SendRequest(req *http.Request, output interface{}) (*http.Response, error) {
	c := http.Client{
		CheckRedirect: func(inreq *http.Request, via []*http.Request) error {
			log.Println("InReq: ", inreq)
			return nil
		}}
	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}
	if resp == nil || resp.Body == nil || output == nil || resp.ContentLength == 0 {
		return resp, err
	}

	contType := resp.Header.Get("Content-Type")
	log.Println("ContType: ", contType)
	defer resp.Body.Close()
	if contType == "application/json" {
		decoder := json.NewDecoder(resp.Body)
		decoder.UseNumber()
		return resp, decoder.Decode(output)
	}
	return resp, nil
}
