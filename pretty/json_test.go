package pretty

import (
	"strings"
	"testing"

	z "github.com/burgesQ/webfmwk/v4/testing"
)

const (
	_simpleJSON        = `{ "content": "json" }`
	_simpleJSONCompact = `{"content":"json"}`
	_simpleJSONPretty  = "{\n  \"content\": \"json\"\n}"

	_nestedJSON        = `{"content": "json", "nested": { "val": 4 } }`
	_nestedJSONCompact = `{"content":"json","nested":{"val":4}}`
	_nestedJSONPretty  = "{\n  \"content\": \"json\",\n  \"nested\": {\n    \"val\": 4\n  }\n}"

	_arrayJSON        = `{ "content": [ "pretty", "json" ]}`
	_arrayJSONCompact = `{"content":["pretty","json"]}`
	_arrayJSONPretty  = "{\n  \"content\": [\n    \"pretty\",\n    \"json\"\n  ]\n}"
)

func TestSimplePrettyJSON(t *testing.T) {
	var tests = map[string]struct {
		input, expected       string
		pretty, errorExpected bool
	}{
		"wrong json format": {
			input:         `{wrong json`,
			errorExpected: true,
		},
		"no json": {
			input:         `no json`,
			errorExpected: true,
		},
		"simple json": {
			input:    _simpleJSON,
			expected: _simpleJSONCompact,
		},
		"simple json pretty": {
			input:    _simpleJSON,
			expected: _simpleJSONPretty,
			pretty:   true,
		},
		"nested json": {
			input:    _nestedJSON,
			expected: _nestedJSONCompact,
		},
		"nested json pretty": {
			input:    _nestedJSON,
			expected: _nestedJSONPretty,
			pretty:   true,
		},

		"array json": {
			input:    _arrayJSON,
			expected: _arrayJSONCompact,
		},
		"array json pretty": {
			input:    _arrayJSON,
			expected: _arrayJSONPretty,
			pretty:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out, e := SimplePrettyJSON(strings.NewReader(test.input), test.pretty)
			z.AssertEqual(t, out, test.expected)
			if test.errorExpected {
				z.AssertNotNil(t, e)
			} else {
				z.AssertNil(t, e)
			}
		})
	}
}
