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
	c := http.Client{}
	resp, err := c.Do(req)
	if resp != nil {
		log.Println("response is: ", resp.Status, err)
		if resp.Body != nil && output == nil {
			defer resp.Body.Close()
		}
	}
	if output != nil {
		return resp, json.NewDecoder(resp.Body).Decode(output)
	}
	return resp, err
}
