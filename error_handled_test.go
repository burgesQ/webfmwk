package webfmwk

import (
	"errors"
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

func TestWrapping(t *testing.T) {
	var (
		testE = errors.New("what a pretty test")
		e     = NewUnauthorized(_testContent).SetWrapped(testE)
		eh    IErrorHandled
	)

	t.Run("test error is", func(t *testing.T) {
		if !errors.Is(e, testE) {
			t.Errorf("ErrorHandled isn't a testE")
		}
	})

	t.Run("test error as", func(t *testing.T) {
		if !errors.As(e, &eh) {
			t.Errorf("Unauthorized isn't an ErrorHandled")
		}
	})

	t.Run("test error unwrap", func(t *testing.T) {
		z.AssertEqual(t, e.Unwrap().Error(), testE.Error())
	})

	// test wrap

	// test Is

	// test As

}

func TestMethod(t *testing.T) {
	var tests = map[string]struct {
		actual, expected int
	}{
		"processing": {
			NewProcessing(_testContent).GetOPCode(), http.StatusProcessing,
		},

		"no contet": {
			NewNoContent().GetOPCode(), http.StatusNoContent,
		},

		"bad request": {
			NewBadRequest(_testContent).GetOPCode(), http.StatusBadRequest,
		},
		"unauthorized": {
			NewUnauthorized(_testContent).GetOPCode(), http.StatusUnauthorized,
		},
		"not found": {
			NewNotFound(_testContent).GetOPCode(), http.StatusNotFound,
		},
		"not acceptable": {
			NewNotAcceptable(_testContent).GetOPCode(), http.StatusNotAcceptable,
		},
		"conflict": {
			NewConflict(_testContent).GetOPCode(), http.StatusConflict,
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
