package webfmwk

import (
	"crypto/tls"
	"testing"

	z "github.com/burgesQ/webfmwk/v4/testing"
)

const _addr = ":4242"

func TestLoadTLS(t *testing.T) {
	// load keys
	s := InitServer()
	worker := s.meta.toServer(_addr)

	// gen ssl keys

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("catched : %v", r)
		}
	}()

	s.loadTLS(&worker, TLSConfig{Key: "./doc/server.key", Cert: "./doc/server.cert"})

	z.AssertSliceEqual(t, worker.TLSConfig.CipherSuites, DefaultCipher)
	z.AssertSliceEqual(t, worker.TLSConfig.CurvePreferences, DefaultCurve)

	z.AssertUInt16Equal(t, worker.TLSConfig.MinVersion, tls.VersionTLS12)
	z.AssertUInt16Equal(t, worker.TLSConfig.MaxVersion, tls.VersionTLS13)
}
