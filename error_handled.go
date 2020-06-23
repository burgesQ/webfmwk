package webfmwk

import (
	"fmt"
	"net/http"
)

type (
	// IErrorHandled interface implement the panic recovering
	ErrorHandled interface {
		// error
		Error() string
		Unwrap() error
		GetOPCode() int
		GetContent() interface{}
		SetWrapped(err error) ErrorHandled
	}

	// ErrorHandled implement the IErrorHandled interface
	errorHandled struct {
		op      int
		content interface{}
		err     error
	}

	// AnonymousError struct is used to answer error
	AnonymousError struct {
		Err string `json:"error"`
		e   error  `json:"-"`
	}

	// Response is returned in case of success
	Response struct {
		Content string `json:"content"`
	}
)

// NewResponse generate a new json response payload
func NewResponse(str string) Response {
	return Response{Content: str}
}

// Error implement the Error interface
func (a AnonymousError) Error() string {
	return a.Err
}

// NewAnonymousError generate a new json error response payload
func NewAnonymousError(err string) AnonymousError {
	return AnonymousError{
		Err: err,
	}
}

// NewAnonymousWrappedError generate a AnonymousError which wrap the err params
func NewAnonymousWrappedError(err error, msg string) AnonymousError {
	return AnonymousError{
		Err: msg,
		e:   err,
	}
}

// NewAnonymousWrappedError generate a AnonymousError which wrap the err params
func NewAnonymousErrorFromError(err error) AnonymousError {
	return AnonymousError{
		Err: err.Error(),
		e:   err,
	}
}

// Error implement the Error interface
func (e errorHandled) Error() string {
	return fmt.Sprintf("[%d]: %#v", e.op, e.content)
}

// Unwrap implemtation the Error interface
func (e errorHandled) Unwrap() error {
	return e.err
}

func (e errorHandled) SetWrapped(err error) ErrorHandled {
	e.err = err
	return e
}

// GetOPCode implement the IErrorHandled interface
func (e errorHandled) GetOPCode() int {
	return e.op
}

// GetContent implement the IErrorHandled interface
func (e errorHandled) GetContent() interface{} {
	return e.content
}

func factory(op int, content interface{}) errorHandled {
	return errorHandled{
		op:      op,
		content: content,
	}
}

// NewError return a new errorHandled var
func NewErrorHandled(op int, content interface{}) ErrorHandled {
	return factory(op, content)
}

// NewProcessing produce an errorHandled with the status code 102
func NewProcessing(content interface{}) ErrorHandled {
	return factory(http.StatusProcessing, content)
}

// NewNoContent produce an errorHandled with the status code 204
func NewNoContent() ErrorHandled {
	return factory(http.StatusNoContent, nil)
}

// NewBadRequest produce an errorHandled with the status code 400
func NewBadRequest(content interface{}) ErrorHandled {
	return factory(http.StatusBadRequest, content)
}

// NewUnauthorized  produce an ErrorHandled with the status code 401
func NewUnauthorized(content interface{}) ErrorHandled {
	return factory(http.StatusUnauthorized, content)
}

// NewNotAcceptable produce an ErrorHandled with the status code 404
func NewNotFound(content interface{}) ErrorHandled {
	return factory(http.StatusNotFound, content)
}

// NewNotAcceptable produce an ErrorHandled with the status code 406
func NewNotAcceptable(content interface{}) ErrorHandled {
	return factory(http.StatusNotAcceptable, content)
}

// NewConflict produce an ErrorHandled with the status code 409
func NewConflict(content interface{}) ErrorHandled {
	return factory(http.StatusConflict, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 422
func NewUnprocessable(content interface{}) ErrorHandled {
	return factory(http.StatusUnprocessableEntity, content)
}

// NewServiceUnavailable produce an ErrorHandled with the status code 422
func NewServiceUnavailable(content interface{}) ErrorHandled {
	return factory(http.StatusServiceUnavailable, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 500
func NewInternal(content interface{}) ErrorHandled {
	return factory(http.StatusInternalServerError, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 501
func NewNotImplemented(content interface{}) ErrorHandled {
	return factory(http.StatusNotImplemented, content)
}
