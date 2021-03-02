package handler

import (
	"strings"

	. "github.com/burgesQ/webfmwk/v4"
)

var (
	ErrMissingContentType = NewNotAcceptable(NewError("Missing Content-Type header"))
	ErrNotJSON            = NewNotAcceptable(NewError("Content-Type is not application/json"))
)

func ContentIsJSON(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		var r = c.GetRequest()

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if ctype := r.Header.Get("Content-Type"); ctype == "" {
				return ErrMissingContentType
			} else if !strings.HasPrefix(ctype, "application/json") {
				return ErrNotJSON
			}
		}

		return next(c)
	})
}
