package webfmwk

import (
	fmtls "crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/burgesQ/gommon/port"
	"github.com/burgesQ/webfmwk/v6/tls"
	"github.com/stretchr/testify/require"
)

func TestSSLEnforced(t *testing.T) {
	s, e := InitServer(
		CheckIsUp(),
		SetPrefix("/api"))

	require.Nil(t, e)

	s.RouteApplier(RoutesPerPrefix{
		"/hello": {
			{Verbe: "GET", Path: "", Handler: func(c Context) error {
				return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
			}},
		},
	})

	defer func() { require.Nil(t, s.ShutdownAndWait()) }()

	port, e := port.GetFree()
	require.Nil(t, e)
	cfg := tls.Config{
		Key:   "./example/ssl.key",
		Cert:  "./example/ssl.crt",
		Ca:    "./example/cacert.pem",
		Level: tls.RequireAndVerifyClientCertAndSAN,
	}

	t.Logf("starting tls server ...\n")

	go s.StartTLS(fmt.Sprintf("127.0.0.1:%d", port), cfg)
	t.Logf("waiting for tls server ...\n")
	<-s.isReady
	t.Logf("tls server is ready ...\n")

	tests := map[string]struct {
		ca        string
		cert, key string
		error     bool
	}{
		"no client cert":  {ca: cfg.GetCa(), error: true},
		"client cert key": {ca: cfg.GetCa(), cert: cfg.GetCert(), key: cfg.GetKey()},
	}

	for name := range tests {
		test := tests[name]

		t.Run(name, func(t *testing.T) {
			caCert, err := os.ReadFile(test.ca)
			require.Nil(t, err)

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			var cert fmtls.Certificate
			if test.cert != "" && test.key != "" {
				cert, err = fmtls.LoadX509KeyPair(test.cert, test.key)
				require.Nil(t, err)
			}

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &fmtls.Config{
						RootCAs:      caCertPool,
						Certificates: []fmtls.Certificate{cert},
					},
				},
			}

			r, e := client.Get(fmt.Sprintf("https://127.0.0.1:%d/api/hello", port))

			if test.error {
				require.NotNil(t, e, "request should fail for "+name)
			} else {
				require.Nil(t, e, "request shouldn't fail for "+name)
				r.Body.Close()
			}
			// todo: test response
		})
	}
}
