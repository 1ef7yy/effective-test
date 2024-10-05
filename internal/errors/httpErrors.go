package errors

type HTTPError struct {
	Code int
	Msg  string
}

func (e HTTPError) Error() string {
	return e.Msg
}
