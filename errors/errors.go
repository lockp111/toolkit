package errors

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// unknow code
const (
	UnknowCode = -1
)

// APIError ...
type APIError struct {
	Code     int32
	Status   int
	Detail   string
	Internal string      `json:",omitempty"`
	Content  interface{} `json:",omitempty"`
}

// ParseAPIError ...
func ParseAPIError(err error) *APIError {
	var (
		retErr *APIError
	)
	switch err.(type) {
	case *APIError:
		retErr = err.(*APIError)
	default:
		retErr = &APIError{
			Code:   UnknowCode,
			Status: 500,
			Detail: err.Error(),
		}
	}
	return retErr
}

func (e *APIError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// Parse ...
func Parse(err string) *APIError {
	api := &APIError{
		Status:   500,
		Code:     UnknowCode,
		Detail:   http.StatusText(500),
		Internal: err,
	}

	if json.Valid([]byte(err)) {
		e := json.Unmarshal([]byte(err), api)
		if e != nil {
			api.Detail = err
		}
	}

	return api
}

// Error ...
type Error struct {
	errCode int32
}

// NewError ...
func NewError(code int32) *Error {
	return &Error{
		errCode: code,
	}
}

// Errors ...
var Errors = map[int32]*APIError{}

func (er *Error) addError(err *APIError) *APIError {
	err.Code += er.errCode
	e, ok := Errors[err.Code]
	if ok {
		log.Fatalf("duplate error code: %v, %v", e, err)
	}

	Errors[err.Code] = err
	return err
}

// BadRequest ...
func (er *Error) BadRequest(code int32, detail string) error {
	return er.addError(&APIError{
		Code:   code,
		Status: 400,
		Detail: detail,
	})
}

// Conflict ..
func (er *Error) Conflict(code int32, detail string) error {
	return er.addError(&APIError{
		Code:   code,
		Status: 409,
		Detail: detail,
	})
}

// Unauthorized ...
func (er *Error) Unauthorized(code int32, detail string) error {
	return er.addError(&APIError{
		Code:   code,
		Status: 401,
		Detail: detail,
	})
}

// Forbidden ..
func (er *Error) Forbidden(code int32, detail string) error {
	return er.addError(&APIError{
		Code:   code,
		Status: 403,
		Detail: detail,
	})
}

// NotFound ..
func (er *Error) NotFound(code int32, detail string) error {
	return er.addError(&APIError{
		Code:   code,
		Status: 404,
		Detail: detail,
	})
}

// Internal ..
func Internal(detail string, err error) error {

	internal := ""
	if err != nil {
		internal = err.Error()
	}

	return &APIError{
		Status:   500,
		Detail:   detail,
		Internal: internal,
	}
}
