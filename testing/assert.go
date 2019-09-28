package testing

import (
	//	"reflect"
	"testing"
)

const (
	notEqual = "assertion failed\ngot :\t>[\t%v\t]<\nwant :\t>[\t%v\t]<"
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
	assertSimpleEqualContext(t, have, want, notEqual, have, want)
}

func assertSimpleNotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleNotEqualContext(t, have, want, notEqual, have, want)
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

// AssertNotEqual run an assertion that the argument are nto equal
func AssertNotEqual(t *testing.T, have, want interface{}) {
	t.Helper()
	assertSimpleNotEqual(t, have, want)
}

// func AssertNotEqualContext(have, want interface{}, context string, args ...interface{}) {
// 	assertSimpleNotEqualContext(have, want, context, args)
// }

// AssertNil run an assertion that the argument is not nil
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

// func AssertTrueContext(have bool, context string, args ...interface{}) {
// 	assertSimpleEqualContext(have, true, context, args...)
// }

// AssertFalse run an assertion that the bool argument is false
func AssertFalse(t *testing.T, have bool) {
	t.Helper()
	assertSimpleEqual(t, have, false)
}

// func AssertFalseContext(have bool, context string, args ...interface{}) {
// 	assertSimpleEqualContext(have, false, context, args...)
// }

// AssertStringEqual run an assertion that the two string arguments are equal
func AssertStringEqual(t *testing.T, have, want string) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// func AssertStringEqualContext(have, want string, context string, args ...interface{}) {
// 	assertSimpleEqualContext(have, want, context, args...)
// }

// AssertStringNotEqual run an assertion that the two string arguments are not equal
func AssertStringNotEqual(t *testing.T, have, want string) {
	t.Helper()
	assertSimpleNotEqual(t, have, want)
}

// func AssertStringNotEqualContext(have, want string, context string, args ...interface{}) {
// 	assertSimpleNotEqualContext(have, want, context, args...)
// }

// AssertStringEqual run an assertion that the two int arguments are equal
func AssertIntEqual(t *testing.T, have, want int) {
	t.Helper()
	assertSimpleEqual(t, have, want)
}

// func AssertIntNotEqual(have, want int) {
// 	assertSimpleNotEqual(have, want)
// }

// func AssertMapEqual(have, want map[string]string) {
// 	Assert(func() bool { return reflect.DeepEqual(have, want) }, notEqual, have, want)
// }
