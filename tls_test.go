package webfmwk

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: start a tls server and assert the server
// TODO: test listener
// TODO: test mTLS settings
func TestLoadTLS(t *testing.T) {
	asserter := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("catched : %v", r)
		}
	}()

	tlsCfg, err := GetTLSCfg(TLSConfig{Key: "./example/server.key", Cert: "./example/server.cert"})

	assert.Nil(t, err)
	assert.Equal(t, tlsCfg.CipherSuites, DefaultCipher)
	asserter.Equal(tlsCfg.CurvePreferences, DefaultCurve)
	asserter.Equal(tlsCfg.MinVersion, uint16(tls.VersionTLS12))
	asserter.Equal(tlsCfg.MaxVersion, uint16(tls.VersionTLS13))
}
