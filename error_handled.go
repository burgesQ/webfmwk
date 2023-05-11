package webfmwk

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	// ErrorHandled interface is used to ease the error processing.
	ErrorHandled interface {
		// Error implelement the Error interface.
		Error() string

		// GetOPCode return the http status code response associated to the error.
		GetOPCode() int

		// SetStatusCode set the error associated http status code.
		SetStatusCode(op int) ErrorHandled

		// GetContent return the error http response content.
		GetContent() interface{}

		// Wrap/Unwrap ?
	}

	handledError struct {
		content interface{}
		op      int
	}

	// Error struct is used to answer http error.
	Error struct {
		e error

		// Message hold the error message.
		//
		// Example: the impossible appened
		Message string `json:"message" example:"no such resource" validate:"required"`

		// Status hold the error code status.
		//
		// Example: 500
		Status int `json:"status" validate:"required"`
	}

	// Response is returned in case of success.
	Response struct {
		// Message hold the error message.
		//
		// Example: action successfully completed
		Message string `json:"content,omitempty"`

		// Status hold the error code status.
		//
		// Example: 200
		Status int `json:"status" example:"204" validate:"required"`
	}
)

// HandleError test if the error argument implement the ErrorHandled interface
// to return a matching response. Otherwise, a 500/internal error is generated
// from the error arguent.
func HandleError(ctx Context, e error) {
	var eh ErrorHandled
	if errors.As(e, &eh) {
		_ = ctx.JSON(eh.GetOPCode(), eh.GetContent())

		return
	}

	_ = ctx.JSONInternalError(NewErrorFromError(e))
}

// NewResponse generate a new Response struct.
func NewResponse(str string) Response {
	return Response{Message: str, Status: http.StatusOK}
}

// SetStatusCode set the response status code.
func (r *Response) SetStatusCode(op int) {
	r.Status = op
}

// Error implement the Error interface.
func (a Error) Error() string {
	return a.Message
}

// NewError generate a Error struct.
func NewError(err string) Error {
	return Error{
		Message: err,
		Status:  http.StatusInternalServerError,
	}
}

// NewCustomWrappedError generate a Error which wrap the err parameter but
// return the msg one.
func NewCustomWrappedError(err error, msg string) Error {
	return Error{
		Message: msg,
		e:       err,
		Status:  http.StatusInternalServerError,
	}
}

// NewErrorFromError generate a Error which wrap the err parameter.
func NewErrorFromError(err error) Error {
	return Error{
		Message: err.Error(),
		e:       err,
		Status:  http.StatusInternalServerError,
	}
}

// SetStatusCode set the error status code.
func (a *Error) SetStatusCode(op int) {
	a.Status = op
}

// Error implement the Error interface.
func (e handledError) Error() string {
	return fmt.Sprintf("[%d]: %#v", e.op, e.content)
}

// SetStatusCode implement the ErrorHandled interface.
func (e handledError) SetStatusCode(op int) ErrorHandled {
	e.op = op

	return e
}

// GetOPCode implement the ErrorHandled interface.
func (e handledError) GetOPCode() int {
	return e.op
}

// GetContent implement the ErrorHandled interface.
func (e handledError) GetContent() interface{} {
	return e.content
}

func factory(op int, content interface{}) handledError {
	ret := handledError{
		op:      op,
		content: content,
	}

	// append status code is possible
	if e, ok := content.(Error); ok {
		e.SetStatusCode(op)
		ret.content = e
	}

	return ret
}

// NewErrorHandled return a struct implementing ErrorHandled with the provided params.
func NewErrorHandled(op int, content interface{}) ErrorHandled {
	return factory(op, content)
}

// NewProcessing produce an ErrorHandled struct with the status code 102.
func NewProcessing(content interface{}) ErrorHandled {
	return factory(http.StatusProcessing, content)
}

// NewNoContent produce an ErrorHandled struct with the status code 204.
func NewNoContent() ErrorHandled {
	return factory(http.StatusNoContent, nil)
}

// NewBadRequest produce an handledError with the status code 400.
func NewBadRequest(content interface{}) ErrorHandled {
	return factory(http.StatusBadRequest, content)
}

// NewUnauthorized  produce an ErrorHandled with the status code 401.
func NewUnauthorized(content interface{}) ErrorHandled {
	return factory(http.StatusUnauthorized, content)
}

// NewForbidden  produce an ErrorHandled with the status code 403.
func NewForbidden(content interface{}) ErrorHandled {
	return factory(http.StatusForbidden, content)
}

// NewNotFound produce an ErrorHandled with the status code 404.
func NewNotFound(content interface{}) ErrorHandled {
	return factory(http.StatusNotFound, content)
}

// NewNotAcceptable produce an ErrorHandled with the status code 406.
func NewNotAcceptable(content interface{}) ErrorHandled {
	return factory(http.StatusNotAcceptable, content)
}

// NewConflict produce an ErrorHandled with the status code 409.
func NewConflict(content interface{}) ErrorHandled {
	return factory(http.StatusConflict, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 422.
func NewUnprocessable(content interface{}) ErrorHandled {
	return factory(http.StatusUnprocessableEntity, content)
}

// NewInternal produce an ErrorHandled with the status code 500.
func NewInternal(content interface{}) ErrorHandled {
	return factory(http.StatusInternalServerError, content)
}

// NewNotImplemented produce an ErrorHandled with the status code 501.
func NewNotImplemented(content interface{}) ErrorHandled {
	return factory(http.StatusNotImplemented, content)
}

// NewServiceUnavailable produce an ErrorHandled with the status code 503.
func NewServiceUnavailable(content interface{}) ErrorHandled {
	return factory(http.StatusServiceUnavailable, content)
}
