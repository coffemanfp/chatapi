package handlers

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter perfoms writing operations for a given http.ResponseWriter.
type ResponseWriter interface {
	// JSON writes a JSON response on the body.
	// 	@param w http.ResponseWriter: Response to write.
	//	@param code int: HTTP status code of the response.
	//	@param v interface{}: object which will be write on the response.
	//	@return err error: writing error.
	JSON(w http.ResponseWriter, code int, v interface{}) error
}

// ResponseWriterImpl is the implementation for the ResponseWriterImpl interface.
type ResponseWriterImpl struct{}

var responseWriterImpl *ResponseWriterImpl

func (rW ResponseWriterImpl) JSON(w http.ResponseWriter, code int, v interface{}) (err error) {
	err = responseJSON(w, code, v)
	return
}

func NewResponseWriterImpl() *ResponseWriterImpl {
	return &ResponseWriterImpl{}
}

// GetResponseWriterImpl gets or initializes a ResponseWriterImpl instance.
func GetResponseWriterImpl() ResponseWriterImpl {
	if responseWriterImpl == nil {
		responseWriterImpl = NewResponseWriterImpl()
	}
	return *responseWriterImpl
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

// Hash is a alias for a map[string]interface{} used to write data objects on the response.
type Hash map[string]interface{}
