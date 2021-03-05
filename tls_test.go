package webfmwk

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: start a tls server and assert the server
// - 1) listen on the correct addr
// - 2) server the correct tls files

func TestLoadTLS(t *testing.T) {
	// load keys
	s := InitServer()
	asserter := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("catched : %v", r)
		}
	}()

	tlsCfg := s.getTLSCfg(TLSConfig{Key: "./example/server.key", Cert: "./example/server.cert"})

	assert.Equal(t, tlsCfg.CipherSuites, DefaultCipher)
	asserter.Equal(tlsCfg.CurvePreferences, DefaultCurve)
	asserter.Equal(tlsCfg.MinVersion, uint16(tls.VersionTLS12))
	asserter.Equal(tlsCfg.MaxVersion, uint16(tls.VersionTLS13))
}
