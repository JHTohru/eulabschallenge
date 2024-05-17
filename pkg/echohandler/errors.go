package echohandler

type malformationError struct {
	msg string
}

func (merr *malformationError) Error() string {
	return merr.msg
}

func newMalformationError(msg string) *malformationError {
	return &malformationError{msg}
}

var (
	errMalformedIDPathParam          = newMalformationError("malformed id path parameter")
	errMalformedRequestBody          = newMalformationError("malformed request body")
	errMalformedPageSizeQueryParam   = newMalformationError("malformed page_size query parameter")
	errMalformedPageNumberQueryParam = newMalformationError("malformed page_number query parameter")
)
