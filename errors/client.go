package errors

import "fmt"

// ClientError represents a error to present to the client.
// Implements the error interface.
type ClientError struct {
	httpCode int
	message  string
}

func (h ClientError) Error() string {
	return h.message
}

// HTTPCode gets the http code of the error. Returns 0 if it's not available.
func (h ClientError) HTTPCode() int {
	return h.httpCode
}

// NewClientError initialices a new error with a ClientError implementation.
//  @param httpCode: represents the http error code for the http response.
//  @param m string: message to be presented to the client.
//  @param a ...interface{}: optional arguments for the message.
//	@return $1 error: new ClientError error implementation instance.
func NewClientError(httpCode int, m string, a ...interface{}) error {
	return ClientError{
		httpCode: httpCode,
		message:  fmt.Sprintf(m, a...),
	}
}
