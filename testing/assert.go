package testing

import (
	"reflect"
	"testing"
)

const (
	notEqual = "failed asserting\nvvvvv\n%v\n^^^^^\nExpected\nvvvvv\n%v\n^^^^^"
)

func Assert(t *testing.T, method func() bool, context string, args ...interface{}) {
	if method() == false {
		t.Errorf(context, args...)
	}
}

func assertSimpleEqualContext(t *testing.T, have, want interface{}, context string, args ...interface{}) {
	Assert(t, func() bool { return have == want }, context, args...)
}

func assertSimpleNotEqualContext(t *testing.T, have, want interface{}, context string, args ...interface{}) {
	Assert(t, func() bool { return have != want }, context, args...)
}

func assertSimpleEqual(t *testing.T, have, want interface{}) {
	assertSimpleEqualContext(t, have, want, notEqual, have, want)
}

func assertSimpleNotEqual(t *testing.T, have, want interface{}) {
	assertSimpleNotEqualContext(t, have, want, notEqual, have, want)
}

func AssertEqual(t *testing.T, have, want interface{}) {
	assertSimpleEqual(t, have, want)
}

func AssertStringEqual(t *testing.T, have, want string) {
	assertSimpleEqual(t, have, want)
}

func AssertStringEqualContext(t *testing.T, have, want string, context string, args ...interface{}) {
	assertSimpleEqualContext(t, have, want, context, args...)
}

func AssertStringNotEqual(t *testing.T, have, want string) {
	assertSimpleNotEqual(t, have, want)
}

func AssertStringNotEqualContext(t *testing.T, have, want string, context string, args ...interface{}) {
	assertSimpleNotEqualContext(t, have, want, context, args...)
}

func AssertIntEqual(t *testing.T, have, want int) {
	assertSimpleEqual(t, have, want)
}

func AssertIntNotEqual(t *testing.T, have, want int) {
	assertSimpleNotEqual(t, have, want)
}

func AssertTrue(t *testing.T, have bool) {
	assertSimpleEqual(t, have, true)
}

func AssertTrueContext(t *testing.T, have bool, context string, args ...interface{}) {
	assertSimpleEqualContext(t, have, true, context, args...)
}

func AssertFalse(t *testing.T, have bool) {
	assertSimpleEqual(t, have, false)
}

func AssertFalseContext(t *testing.T, have bool, context string, args ...interface{}) {
	assertSimpleEqualContext(t, have, false, context, args...)
}

func AssertMapEqual(t *testing.T, have, want map[string]string) {
	Assert(t, func() bool { return reflect.DeepEqual(have, want) }, notEqual, have, want)
}
