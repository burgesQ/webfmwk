package webfmwk

import "net/http"

type (
	// Interface ErrorHandled implement the panic recovering
	IErrorHandled interface {
		GetOPCode() int
		GetContent() interface{}
	}

	// Error implement the ErrorHandled interface
	ErrorHandled struct {
		op      int
		content interface{}
	}
)

// GetOPCode implement the IErrorHandled interface
func (e ErrorHandled) GetOPCode() int {
	return e.op
}

// GetContentimplement the IErrorHandled interface
func (e ErrorHandled) GetContent() interface{} {
	return e.content
}

func factory(op int, content interface{}) ErrorHandled {
	return ErrorHandled{
		op:      op,
		content: content,
	}
}

// NewNoContent produce an ErrorHandled with the status code 204
func NewNoContent() ErrorHandled {
	return factory(http.StatusNoContent, nil)
}

// NewBadRequest produce an ErrorHandled with the status code 400
func NewBadRequest(content interface{}) ErrorHandled {
	return factory(http.StatusBadRequest, content)
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
