package webfmwk

import (
	"net/http"
	"testing"

	z "github.com/burgesQ/webfmwk/v2/testing"
)

var (
	_testOP              = 200
	_testContent         = "ok"
	_testingErrorHandled = ErrorHandled{
		op:      _testOP,
		content: _testContent,
	}
)

func TestGetOPCode(t *testing.T) {
	z.AssertEqual(t, _testOP, _testingErrorHandled.GetOPCode())
}

func TestGetContent(t *testing.T) {
	z.AssertEqual(t, _testContent, _testingErrorHandled.GetContent())
}

func TestFactory(t *testing.T) {
	test := factory(_testOP, _testContent)
	z.AssertEqual(t, _testOP, test.GetOPCode())
	z.AssertEqual(t, _testContent, test.GetContent())
}

func TestMethod(t *testing.T) {
	var tests = map[string]struct {
		actual, expected int
	}{
		"no contet": {
			NewNoContent().GetOPCode(), http.StatusNoContent,
		},
		"bad request": {
			NewBadRequest(_testContent).GetOPCode(), http.StatusBadRequest,
		},
		"not found": {
			NewNotFound(_testContent).GetOPCode(), http.StatusNotFound,
		},

		"not acceptable": {
			NewNotAcceptable(_testContent).GetOPCode(), http.StatusNotAcceptable,
		},
		"unprocessable": {
			NewUnprocessable(_testContent).GetOPCode(), http.StatusUnprocessableEntity,
		},

		"internal": {
			NewInternal(_testContent).GetOPCode(), http.StatusInternalServerError,
		},
		"new implemented": {
			NewNotImplemented(_testContent).GetOPCode(), http.StatusNotImplemented,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			z.AssertEqual(t, test.actual, test.expected)
		})
	}
}
