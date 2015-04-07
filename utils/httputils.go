package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendJsonResponse(rw http.ResponseWriter, value interface{}) {
	js, err := json.Marshal(value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func SendJsonError(rw http.ResponseWriter, value interface{}, code int) {
	js, err := json.Marshal(value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	log.Println("Json: ", string(js))
	http.Error(rw, string(js), code)
}
