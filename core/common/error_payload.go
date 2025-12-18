package common

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func (p ErrorPayload) Error() string {
	if p.Message != "" {
		return p.Message
	}

	if p.Code != "" {
		return p.Code
	}

	return "unknown error"
}
func (p ErrorPayload) WithCause(cause error) ErrorPayload {
	p.Cause = cause
	return p
}
func (p ErrorPayload) WithCode(code string) ErrorPayload {
	p.Code = code
	return p
}
func (p ErrorPayload) WithMessage(message string) ErrorPayload {
	p.Message = message
	return p
}

var _ error = ErrorPayload{}
