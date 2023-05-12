package tls

import (
	"crypto/tls"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	var lvl Level
	e := json.Unmarshal([]byte(`"allow"`), &lvl)
	require.Nil(t, e)
	require.Equal(t, RequireAnyClientCert, lvl)

	b, e := json.Marshal(VerifyClientCertIfGiven)
	require.Nil(t, e)
	require.Equal(t, `"try"`, string(b))

	require.Equal(t, tls.RequireAnyClientCert, lvl.STD())
}
