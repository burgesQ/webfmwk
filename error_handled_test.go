package webfmwk

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, _testOP, _testingErrorHandled.GetOPCode())
}

func TestGetContent(t *testing.T) {
	assert.Equal(t, _testContent, _testingErrorHandled.GetContent())
}

func TestFactory(t *testing.T) {
	test := factory(_testOP, _testContent)
	assert.Equal(t, _testOP, test.GetOPCode())
	assert.Equal(t, _testContent, test.GetContent())
}

func TestNewErrorHandled(t *testing.T) {
	e := NewErrorHandled(_testOP, _testContent)
	assert.Equal(t, _testOP, e.GetOPCode())
	assert.Equal(t, _testContent, e.GetContent())
	assert.Equal(t, `[200]: "ok"`, e.Error())
}

func TestResponse(t *testing.T) {
	assert.Equal(t, "test", NewResponse("test").Message)
}

func TestError(t *testing.T) {
	var err = errors.New("test")

	asserter := assert.New(t)

	e := NewError("testing")
	asserter.True(e.Message == "testing")
	e = NewCustomWrappedError(err, "testing")
	asserter.True(e.e == err)
	asserter.Equal("testing", e.Message)
	e = NewErrorFromError(err)
	asserter.Equal("test", e.Message)
	asserter.True(e.e == err)
	asserter.Equal("test", e.Message)
	asserter.Equal("test", e.Error())
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
		"forbidden": {
			NewForbidden(_testContent).GetOPCode(), http.StatusForbidden,
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
		"not implemented": {
			NewNotImplemented(_testContent).GetOPCode(), http.StatusNotImplemented,
		},
		"service unavailable": {
			NewServiceUnavailable(_testContent).GetOPCode(), http.StatusServiceUnavailable,
		},
	}

	for name, te := range tests {
		test := te
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.actual)
		})
	}
}
