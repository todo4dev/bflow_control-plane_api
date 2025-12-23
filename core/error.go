package core

type Error struct {
	Code    string
	Message string
	cause   error
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return "unknown error"
}
func (e Error) Unwrap() error {
	return e.cause
}
func (e Error) WithCause(error error) Error {
	e.cause = error
	return e
}
func (e Error) WithCode(code string) Error {
	e.Code = code
	return e
}
func (e Error) WithMessage(message string) Error {
	e.Message = message
	return e
}
