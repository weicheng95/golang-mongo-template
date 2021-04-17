package httperr

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

type HttpError struct {
	error     string
	message   string
	errorType ErrorType
}

func (s HttpError) Error() string {
	return s.error
}

func (s HttpError) Message() string {
	return s.message
}

func (s HttpError) ErrorType() ErrorType {
	return s.errorType
}