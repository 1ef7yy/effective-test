package errors

import "fmt"

type APIError interface {
	Error() string
	Status() int
	Message() string
}

type HTTPError struct {
	Code int
	Msg  string
}

func (e HTTPError) Error() string {
	return fmt.Sprint("code: ", e.Code, " msg: ", e.Msg)
}

func (e HTTPError) Status() int {
	return e.Code

}

func (e HTTPError) Message() string {
	return e.Msg
}

func NewHTTPError(code int, msg string) HTTPError {
	return HTTPError{
		Code: code,
		Msg:  msg,
	}
}
