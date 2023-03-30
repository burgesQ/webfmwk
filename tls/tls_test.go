package tls

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
		icfg := Config{
			Key:      "../example/server.key",
			Cert:     "../example/server.cert",
			Ca:       "../example/cacert.pem",
			Insecure: true,
		}
		cfg, err := GetTLSCfg(icfg)

		asserter.Nil(err)
		asserter.Equal(DefaultCipher, cfg.CipherSuites)
		asserter.Equal(DefaultCurve, cfg.CurvePreferences)
		asserter.Equal(uint16(tls.VersionTLS12), cfg.MinVersion)
		asserter.Equal(uint16(tls.VersionTLS13), cfg.MaxVersion)
		// TODO: test loaded certs ?
		asserter.Equal(tls.NoClientCert, cfg.ClientAuth)
		asserter.Equal("\t ~!~ cert:\t\"../example/server.cert\"\n\t ~!~ key:\t\"../example/server.key\"\n"+
			"\t ~!~ ca:\t\"../example/cacert.pem\",\n\t ~!~ insecure:\ttrue\n\t ~!~ level:\tnever\n", icfg.String())
	}

	/* TODO */
	// t.Log("secured mTLS config")	{}
}
