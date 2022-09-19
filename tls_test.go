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

	t.Log("insecure config")
	{
		tlsCfg, err := GetTLSCfg(TLSConfig{
			Key:      "./example/server.key",
			Cert:     "./example/server.cert",
			Ca:       "./example/cacert.pem",
			Insecure: true,
		})

		asserter.Nil(err)
		asserter.Equal(DefaultCipher, tlsCfg.CipherSuites)
		asserter.Equal(DefaultCurve, tlsCfg.CurvePreferences)
		asserter.Equal(uint16(tls.VersionTLS12), tlsCfg.MinVersion)
		asserter.Equal(uint16(tls.VersionTLS13), tlsCfg.MaxVersion)
		// TODO: test loaded certs ?
		asserter.Equal(tlsCfg.ClientAuth, tls.NoClientCert)
	}

	t.Log("secured mTLS config")
	{
		/* TODO */
	}
}
