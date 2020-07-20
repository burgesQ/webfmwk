package webfmwk

import (
	"errors"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/assert"
)

var (
	_testOP              = 200
	_testContent         = "ok"
	_testingErrorHandled = errorHandled{
		op:      _testOP,
		content: _testContent,
	}
)

func TestGetOPCode(t *testing.T) {
	assert.Equal(t, _testingErrorHandled.GetOPCode(), _testOP)
}

func TestGetContent(t *testing.T) {
	assert.Equal(t, _testingErrorHandled.GetContent(), _testContent)
}

func TestFactory(t *testing.T) {
	test := factory(_testOP, _testContent)
	assert.Equal(t, test.GetOPCode(), _testOP)
	assert.Equal(t, test.GetContent(), _testContent)
}

func TestNewErrorHandled(t *testing.T) {
	e := NewErrorHandled(_testOP, _testContent)
	assert.Equal(t, e.GetOPCode(), _testOP)
	assert.Equal(t, e.GetContent(), _testContent)
	assert.Equal(t, e.Error(), `[200]: "ok"`)
}

func TestWrapping(t *testing.T) {
	var (
		testE = errors.New("what a pretty test")
		e     = NewUnauthorized(_testContent).SetWrapped(testE)
		eh    ErrorHandled
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
		assert.Equal(t, e.Unwrap().Error(), testE.Error())
	})

	// test wrap

	// test Is

	// test As

}

func TestResponse(t *testing.T) {
	assert.StringEqual(t, NewResponse("test").Message, "test")
}

func TestError(t *testing.T) {
	var err = errors.New("test")

	e := NewError("testing")
	assert.True(t, e.Message == "testing")
	e = NewAnonymousWrappedError(err, "testing")
	assert.True(t, e.e == err)
	assert.StringEqual(t, e.Message, "testing")
	e = NewErrorFromError(err)
	assert.StringEqual(t, e.Message, "test")
	assert.True(t, e.e == err)
	assert.StringEqual(t, e.Message, "test")
	assert.StringEqual(t, e.Error(), "test")
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
		"service unavailable": {
			NewServiceUnavailable(_testContent).GetOPCode(), http.StatusServiceUnavailable,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.actual, test.expected)
		})
	}
}
