package webfmwktest

import (
	"testing"

	"github.com/burgesQ/webfmwk/v4"
)

func TestBasic(t *testing.T) {
	testHandler := func(c webfmwk.Context) error {
		return c.JSONOk(webfmwk.Response{Message: "ok"})
	}

	// not test handler but : -->
	//	CustomHandler(handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	// create context (s.CustomHandler)
	// }

	GetAndTest(t, testHandler, Expected{
		Code: 200,
		Body: `{"status":0,"content":"ok"}`,
	})
}
