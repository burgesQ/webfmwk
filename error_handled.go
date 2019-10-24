package webfmwk

import (
	"fmt"
	"net/http"
)

type (
	// IErrorHandled interface implement the panic recovering
	IErrorHandled interface {
		GetOPCode() int
		GetContent() interface{}
	}

	// ErrorHandled implement the IErrorHandled interface
	ErrorHandled struct {
		op      int
		content interface{}
	}
)

// GetOPCode implement the IErrorHandled interface
func (e ErrorHandled) GetOPCode() int {
	return e.op
}

// GetContent implement the IErrorHandled interface
func (e ErrorHandled) GetContent() interface{} {
	return e.content
}

// func (e ErrorHandled) String() string {
// 	if utf8.Valid(e.content) {
// 		return fmt.Sprintf("[%d]: %s", e.op, e.content)
// 	}
// 	return fmt.Sprintf("[%d]: %#v", e.op, e.content)
// }

func (e ErrorHandled) Error() string {
	return fmt.Sprintf("[%d]: %#v", e.op, e.content)
}

func factory(op int, content interface{}) ErrorHandled {
	return ErrorHandled{
		op:      op,
		content: content,
	}
}

// NewError return a new ErrorHandled var
func NewErrorHandled(op int, content interface{}) ErrorHandled {
	return factory(op, content)
}

func NewProcessing(content interface{}) ErrorHandled {
	return factory(http.StatusProcessing, content)
}

// NewNoContent produce an ErrorHandled with the status code 204
func NewNoContent() ErrorHandled {
	return factory(http.StatusNoContent, nil)
}

// NewBadRequest produce an ErrorHandled with the status code 400
func NewBadRequest(content interface{}) ErrorHandled {
	return factory(http.StatusBadRequest, content)
}

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

// NewUnprocessable produce an ErrorHandled with the status code 422
func NewUnprocessable(content interface{}) ErrorHandled {
	return factory(http.StatusUnprocessableEntity, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 500
func NewInternal(content interface{}) ErrorHandled {
	return factory(http.StatusInternalServerError, content)
}

// NewUnprocessable produce an ErrorHandled with the status code 501
func NewNotImplemented(content interface{}) ErrorHandled {
	return factory(http.StatusNotImplemented, content)
}
