package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	sErrors "github.com/coffemanfp/chat/errors"
)

// RequestReader perfoms reading operations for a given *http.Request.
type RequestReader interface {
	// JSON reads the body information as a JSON, writes output on the v interface{} param.
	// 	@param r *http.Request: Request to read.
	//	@param v interface{}: output target. (Must be a pointer)
	//	@return err error: reading error.
	JSON(r *http.Request, v interface{}) error
}

// RequestReaderImpl is the implementation for the RequestReader interface.
type RequestReaderImpl struct{}

var requestReaderImpl *RequestReaderImpl

func (rR RequestReaderImpl) JSON(r *http.Request, v interface{}) (err error) {
	if r == nil {
		err = fmt.Errorf("invalid request value: empty or nil *http.Request")
		err = sErrors.NewClientError(http.StatusInternalServerError, sErrors.SERVER_ERROR_MESSAGE, err)
		return
	}
	if !checkContentTypeJSON(r.Header) {
		err = fmt.Errorf("invalid content type: Content-Type header is not application/json")
		return
	}

	err = json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("error checking body: empty body content (%s)", err)
			return
		}

		err = fmt.Errorf("error decoding body: %s", err)
	}
	return
}

// NewRequestReaderImpl initializes a new RequestReaderImpl instance.
// 	@return $1 RequestReaderImpl: new RequestReaderImpl instance.
func NewRequestReaderImpl() *RequestReaderImpl {
	return &RequestReaderImpl{}
}

// GetRequestReaderImpl gets or initializes a RequestReaderImpl instance.
func GetRequestReaderImpl() RequestReaderImpl {
	if requestReaderImpl == nil {
		requestReaderImpl = NewRequestReaderImpl()
	}
	return *requestReaderImpl
}

func checkContentTypeJSON(h http.Header) (match bool) {
	return h.Get("Content-Type") == "application/json"
}
