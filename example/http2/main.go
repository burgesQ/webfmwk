package main

/**
 * This is a small poc to test mTLS w/ http2 via the fasthttp go' implementation.
 *
 * Started on: Thu Jul 4 9:30
 * Ended on: Wed Jul 5 11:40
 *
 * Usefull links
 * - extra tls config nextProto
 *   https://github.com/dgrr/http2/blob/master/configure.go#L111
 * - not working example
 *   https://github.com/dgrr/http2/issues/24
 * - curl and http2
 *   https://curl.se/docs/http2.html
 * - curl pretty header & payload
 *   http://blog.aaronholmes.net/displaying-response-headers-and-pretty-json-with-curl/
 * - curl TLS and mTLS
 *   https://smallstep.com/hello-mtls/doc/client/curl
 *
 * The following needed extra attention:
 * - tls config need next proto h2
 * - mtls wrapper (GetConfigForClient) also needed next proto h2
 * - listenTLS didn't work w/ custom tls config
 * - tls listner didn't work w/ listen if the listner was created before the h2 setup
 *
 * See ya next time for more fun ride.
 *
 *
 * Q
 **/

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"

	fasthttp2 "github.com/dgrr/http2"
	"github.com/valyala/fasthttp"
)

func main() {
	tlsCfg := loadTLS()

	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, "Hello, world!\n\n")

			ctx.SetContentType("text/plain; charset=utf8")
		},
	}
	// fasthttp2.ConfigureServerAndConfig(server, &tlsCfg)
	fasthttp2.ConfigureServer(server, fasthttp2.ServerConfig{Debug: true})

	// fmt.Printf("\n%+v\n", tlsCfg)

	addr := "192.168.56.1:4443"

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("listing on https://%q\n", addr)

	tlsListener := tls.NewListener(listener, tlsCfg)

	if err := server.Serve(tlsListener); err != nil {
		panic(err)
	}
}

func loadTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair("/home/master/repo/certs/ssl.crt",
		"/home/master/repo/certs/ssl.key")
	if err != nil {
		panic(fmt.Errorf("cannot load cert and key: %w", err))
	}

	tlsCfg := &tls.Config{
		NextProtos:   []string{"h2"},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}

	return loadMTLSVerifyPeer(loadCA(tlsCfg))
}

func loadMTLSVerifyPeer(tlsCfg *tls.Config) *tls.Config {
	tlsCfg.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		return &tls.Config{
			NextProtos: []string{"h2"},
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		}, nil
	}

	return tlsCfg
}

func loadCA(tlsCfg *tls.Config) *tls.Config {
	pool := x509.NewCertPool()

	if caCertPEM, e := os.ReadFile("/home/master/repo/certs/cacert.pem"); e != nil {
		panic(fmt.Errorf("cannot load ca cert in pool: %w", e))
	} else if !pool.AppendCertsFromPEM(caCertPEM) {
		panic(errors.New("cannot apend cert to pem"))
	}

	tlsCfg.ClientCAs = pool

	return tlsCfg
}
