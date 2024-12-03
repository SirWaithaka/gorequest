package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Error reports an error or the api response and the status code of the request
type Error struct {
	statusCode int
	response   []byte
	Err        error
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %v", e.statusCode, string(e.response), e.Err)
}

//// Response converts http responses into ErrorResponse type.
//func (e Error) Response() ErrorResponse {
//	if e.response == nil {
//		return ErrorResponse{}
//	}
//
//	var res ErrorResponse
//	if err := jsoniter.NewDecoder(bytes.NewReader(e.response)).Decode(&res); err != nil {
//		return ErrorResponse{}
//	}
//	return res
//}

func (e Error) Timeout() bool {
	// check if status code is gateway timeout
	if e.statusCode == http.StatusGatewayTimeout {
		return true
	}

	// check if error is a timeout
	var err *url.Error
	return errors.As(e.Err, &err) && err.Timeout()
}
