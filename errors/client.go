package errors

import "fmt"

type ClientError struct {
	httpCode int
	message  string
}

func (h ClientError) Error() string {
	return h.message
}

func (h ClientError) HTTPCode() int {
	return h.httpCode
}

func NewClientError(httpCode int, m string, a ...interface{}) error {
	return ClientError{
		httpCode: httpCode,
		message:  fmt.Sprintf(m, a...),
	}
}
