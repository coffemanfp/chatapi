package handlers

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter interface {
	JSON(w http.ResponseWriter, v interface{}) error
}

type ResponseWriterImpl struct{}

func (rW ResponseWriterImpl) JSON(w http.ResponseWriter, v interface{}) (err error) {
	err = responseJSON(w, v)
	return
}

func NewResponseWriterImpl() ResponseWriterImpl {
	return ResponseWriterImpl{}
}

func responseJSON(w http.ResponseWriter, v interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	raw, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return
	}
	_, err = w.Write(raw)
	return
}

type Hash map[string]interface{}
