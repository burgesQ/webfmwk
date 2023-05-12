package cmd

import (
	"reflect"
	"testing"

	"github.com/burgesQ/webfmwk/v5/tls"
	"github.com/stretchr/testify/require"
)

func TestStringToLevelHookFunc(t *testing.T) {
	// a := reflect.TypeOf("ca_cert")
	// b := reflect.ValueOf(tls.RequestClientCert)

	t.Log("classic processing")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(tls.RequireAnyClientCert),
			"hardAndSAN")
		require.Nil(t, e)
		require.Equal(t, tls.RequireAndVerifyClientCertAndSAN, l)
	}

	t.Log("failed processing (1)")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf(1),
			reflect.TypeOf(2),
			-3)
		require.Nil(t, e)
		require.Equal(t, -3, l)
	}

	t.Log("failed processing (2)")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(tls.RequireAnyClientCert),
			-3)
		require.Nil(t, e)
		require.Equal(t, -3, l)
		// t.Logf("%+v", l)
	}

	t.Log("failed processing (3)")
	{
		_, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(tls.RequireAnyClientCert),
			"abcd")
		require.NotNil(t, e)
		// require.Equal(t, "abcd", l)
		// t.Logf("%+v", l)
	}

	t.Log("processing skippped because of no verif")
	{
		l, e := stringToLevelHookFunc(
			reflect.TypeOf("level"),
			reflect.TypeOf(tls.RequireAnyClientCert),
			tls.NoClientCert)
		require.Nil(t, e)
		require.Equal(t, tls.NoClientCert, l)
	}
}
