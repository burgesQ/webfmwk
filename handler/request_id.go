package handler

import (
	webfmwk "github.com/burgesQ/webfmwk/v3"
	"github.com/google/uuid"
)

func RequestID(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.IContext) {
		c.SetRequestID(uuid.New().String())
		next(c)
	})
}
