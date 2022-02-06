package handlers

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter interface {
	JSON(w http.ResponseWriter, code int, v interface{}) error
}

type ResponseWriterImpl struct{}

func (rW ResponseWriterImpl) JSON(w http.ResponseWriter, code int, v interface{}) (err error) {
	err = responseJSON(w, code, v)
	return
}

func NewResponseWriterImpl() ResponseWriterImpl {
	return ResponseWriterImpl{}
}

func responseJSON(w http.ResponseWriter, code int, v interface{}) (err error) {
	w.Header().Add("Content-Type", "application/json")
	raw, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(raw)
	return
}

type Hash map[string]interface{}
