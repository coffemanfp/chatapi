package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RequestReader interface {
	JSON(r *http.Request, v interface{}) error
}

type RequestReaderImpl struct{}

func (rR RequestReaderImpl) JSON(r *http.Request, v interface{}) (err error) {
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

func NewRequestReaderImpl() RequestReaderImpl {
	return RequestReaderImpl{}
}

func checkContentTypeJSON(h http.Header) (match bool) {
	fmt.Printf("Content-Type is %s\n", h.Get("Content-Type"))
	return h.Get("Content-Type") == "application/json"
}
