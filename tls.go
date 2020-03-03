package webfmwk

import (
	"crypto/tls"
	"net/http"
)

type (
	// TLSConfig contain the tls config passed by the config file
	TLSConfig struct {
		Cert     string `json:"cert"`
		Key      string `json:"key"`
		Insecure bool   `json:"insecure"`
		// CaCert string `json:"ca-cert"`
	}
)

// StartTLS expose an server to an HTTPS endpoint
func (s *Server) StartTLS(addr string, tlsStuffs TLSConfig) {
	s.internalHandler()
	s.launcher.Start("https server "+addr, func() error {
		go s.pollPingEndpoint(addr)
		return s.internalInit(addr, tlsStuffs).ListenAndServeTLS(tlsStuffs.Cert, tlsStuffs.Key)
	})
}

func (s *Server) loadTLS(worker *http.Server, tlsCfg TLSConfig) {
	/* #nosec */
	worker.TLSConfig = &tls.Config{
		InsecureSkipVerify:       tlsCfg.Insecure,
		Certificates:             make([]tls.Certificate, 1),
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
	}

	cert, err := tls.LoadX509KeyPair(tlsCfg.Cert, tlsCfg.Key)
	if err != nil {
		s.log.Fatalf("%s", err.Error())
	}

	worker.TLSConfig.Certificates[0] = cert
}
