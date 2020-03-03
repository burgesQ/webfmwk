package testing

import (
	"reflect"
	"testing"
)

const (
	STRING = 1
	INT    = 2

	_notEqual = "assertion failed\ngot :\t>[\t%v\t]<\nwant :\t>[\t%v\t]<"
	_flagType = "STRING = 1 - INT = 2"
)

func assert(t *testing.T, method func() bool,
	context string, args ...interface{}) {
	t.Helper()
	if !method() {
		if len(args) > 0 {
			t.Errorf(context, args...)
		} else {
			t.Errorf(context)
		}
	}
}

func assertSimpleEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have == want }, context, args...)
}

func assertSimpleNotEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assert(t, func() bool { return have != want }, context, args...)
}

func assertSimpleEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleEqualContext(t, have, want, _notEqual, have, want)
}

func assertSimpleNotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleNotEqualContext(t, have, want, _notEqual, have, want)
}

// AssertEqual run an assertion that the argument are equal
func AssertEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// AssertEqualContext run a assertion that the arguments are equal
// on fail the error message context is yield
func AssertEqualContext(t *testing.T, have, want interface{},
	context string, args ...interface{}) {
	t.Helper()
	assertSimpleEqualContext(t, have, want, context, args...)
}

// AssertNotEqual run an assertion that the argument are not equal
func AssertNotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleNotEqual(t, have, want)
}

// AssertNotNil run an assertion that the argument is not nil
func AssertNotNil(t *testing.T, have interface{}) {
	t.Helper()
	assertSimpleNotEqual(t, have, nil)
}

// AssertNil run an assertion that the argument is nil
func AssertNil(t *testing.T, have interface{}) {
	t.Helper()
	assertSimpleEqual(t, have, nil)
}

// AssertTrue run an assertion that the bool argument is true
func AssertTrue(t *testing.T, have bool) {
	t.Helper()
	assertSimpleEqual(t, have, true)
}

// AssertTrueContext run an assertion that the bool argument is true
func AssertTrueContext(t *testing.T, have bool, fmt string, args ...interface{}) {
	t.Helper()
	assertSimpleEqualContext(t, have, true, fmt, args...)
}

// AssertFalse run an assertion that the bool argument is false
func AssertFalse(t *testing.T, have bool) {
	t.Helper()
	assertSimpleEqual(t, have, false)
}

// AssertFalseContext run an assertion that the bool argument is false with the
// provied context
func AssertFalseContext(t *testing.T, have bool, context string, args ...interface{}) {
	t.Helper()
	assertSimpleEqualContext(t, have, false, context, args...)
}

// AssertStringEqual run an assertion that the two string arguments are equal
func AssertStringEqual(t *testing.T, have, want string) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// AssertStringNotEqual run an assertion that the two string arguments are not equal
func AssertStringNotEqual(t *testing.T, have, want string) {
	t.Helper()
	assertSimpleNotEqual(t, have, want)
}

// AssertIntEqual run an assertion that the two string arguments are equal
func AssertIntEqual(t *testing.T, have, want int) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// AssertUInt16Equal run an assertion that the two string arguments are equal
func AssertUInt16Equal(t *testing.T, have, want uint16) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// AssertIntNotEqual run an assertion that the two string arguments are not equal
func AssertIntNotEqual(t *testing.T, have, want int) {
	t.Helper()
	assertSimpleNotEqual(t, have, want)
}

// AssertMapOfInterfaceEqual
func AssertMapOfInterfaceEqual(t *testing.T, have, want map[string]interface{}) {
	t.Helper()
	AssertTrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// AssertMapOfStringEqual
func AssertMapOfStringEqual(t *testing.T, have, want map[string]string) {
	t.Helper()
	AssertTrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// AssertIs assert the type of have
func AssertIs(t *testing.T, have interface{}, what int) {
	t.Helper()
	var (
		fmt = ""
		ok  = false
	)
	switch what {
	case STRING:
		_, ok = have.(string)
		fmt = "%s is not a string (%T)"
	case INT:
		_, ok = have.(int)
		fmt = "%s is not a int (%T)"
	default:
		t.Errorf("type assertion not supported (%d)[%s]", what, _flagType)
	}
	AssertTrueContext(t, ok, fmt, have, have)
}

// AssertIsString assert that have is a string
func AssertIsString(t *testing.T, have interface{}) {
	t.Helper()
	AssertIs(t, have, STRING)
}

// AssertIsInt assert that have in as int
func AssertIsInt(t *testing.T, have interface{}) {
	t.Helper()
	AssertIs(t, have, INT)
}

// 0 - 0 -- false
// 1 - 0 -- false
// 0 - 1 -- true
func AssertLower(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have < want }, "%d is greater than %d ", have, want)
}

// 0 - 0 -- true
// 1 - 0 -- true
// 0 - 1 -- false
func AssertEqualOrGreater(t *testing.T, have, want int) {
	t.Helper()
	assert(t, func() bool { return have >= want }, "%d is strictly lower than %d ", have, want)
}

// AssertSliceByteEqual
func AssertSliceByteEqual(t *testing.T, have, want []byte) {
	t.Helper()
	AssertTrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}

// AssertSliceEqual
func AssertSliceEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	AssertTrueContext(t, func() bool {
		return reflect.DeepEqual(have, want)
	}(), _notEqual, have, want)
}
