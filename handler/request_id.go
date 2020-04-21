package handler

import (
	webfmwk "github.com/burgesQ/webfmwk/v4"
	"github.com/google/uuid"
)

// RequestID add a unique identifier to the context object
func RequestID(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		c.SetRequestID(uuid.New().String())
		return next(c)
	})
}
