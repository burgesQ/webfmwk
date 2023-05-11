package tls

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: start a tls server and require the server
// TODO: test listener
// TODO: test mTLS settings
func TestLoadTLS(t *testing.T) {
	requirer := require.New(t)

	t.Log("insecure config")
	{
		icfg := Config{
			Key:      "../example/server.key",
			Cert:     "../example/server.cert",
			Insecure: true,
		}
		cfg, err := GetTLSCfg(icfg)

		requirer.Nil(err)
		requirer.Equal(DefaultCipher, cfg.CipherSuites)
		requirer.Equal(DefaultCurve, cfg.CurvePreferences)
		requirer.Equal(uint16(tls.VersionTLS12), cfg.MinVersion)
		requirer.Equal(uint16(tls.VersionTLS13), cfg.MaxVersion)
		// TODO: test loaded certs ?
		requirer.Equal(tls.NoClientCert, cfg.ClientAuth)
		requirer.Equal("\t ~!~ cert:\t\"../example/server.cert\"\n\t ~!~ key:\t\"../example/server.key\"\n"+
			"\t ~!~ ca:\t\"\",\n\t ~!~ insecure:\ttrue\n\t ~!~ level:\tnever\n", icfg.String())
	}

	t.Log("secure config")
	{
		icfg := Config{
			Key:  "../example/ssl.key",
			Cert: "../example/ssl.crt",
			Ca:   "../example/cacert.pem",
			// Insecure: true,
			Level: RequireAndVerifyClientCertAndSAN,
		}
		cfg, err := GetTLSCfg(icfg)

		requirer.Nil(err)
		requirer.Equal(DefaultCipher, cfg.CipherSuites)
		requirer.Equal(DefaultCurve, cfg.CurvePreferences)
		requirer.Equal(uint16(tls.VersionTLS12), cfg.MinVersion)
		requirer.Equal(uint16(tls.VersionTLS13), cfg.MaxVersion)
		// TODO: test loaded certs ?
		requirer.Equal(tls.RequireAndVerifyClientCert, cfg.ClientAuth)
		requirer.Equal("\t ~!~ cert:\t\"../example/ssl.crt\"\n\t ~!~ key:\t\"../example/ssl.key\"\n"+
			"\t ~!~ ca:\t\"../example/cacert.pem\",\n\t ~!~ insecure:\tfalse\n\t ~!~ level:\thardAndSAN\n", icfg.String())
	}

	/* TODO */
	// t.Log("secured mTLS config")	{}
}
