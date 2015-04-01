package utils

import (
	"encoding/json"
	"net/http"
)

func SendJsonResponse(rw http.ResponseWriter, value map[string]interface{}) {
	js, err := json.Marshal(value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
