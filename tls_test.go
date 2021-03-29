package webfmwk

import (
	"crypto/tls"
	"testing"

	"github.com/burgesQ/gommon/assert"
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

	s.loadTLS(&worker, TLSConfig{
		Key:      "./example/server.key",
		Cert:     "./example/server.cert",
		Insecure: true,
		// TODO: cert ca cert pool
	})

	assert.SliceU16Equal(t, worker.TLSConfig.CipherSuites, DefaultCipher)
	//	assert.SliceU16Equal(t, worker.TLSConfig.CurvePreferences, DefaultCurve)

	assert.UInt16Equal(t, worker.TLSConfig.MinVersion, tls.VersionTLS12)
	assert.UInt16Equal(t, worker.TLSConfig.MaxVersion, tls.VersionTLS13)
}
