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

var (
	DefaultCurve = []tls.CurveID{
		tls.CurveP256,
		tls.X25519,
	}
	DefaultCipher = []uint16{

		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,   // HTTP/2-required AES_128_GCM_SHA256 cipher
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   // ECDHE-RSA-AES256-GCM-SHA384
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   // ECDH-RSA-AES256-GCM-SHA384
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, // ECDHE-ECDSA-AES256-GCM-SHA384
		tls.TLS_AES_128_GCM_SHA256,                  // 1.3 tls cipher
		tls.TLS_AES_256_GCM_SHA384,                  // 1.3 tls cipher
		tls.TLS_CHACHA20_POLY1305_SHA256,            // 1.3 tls cipher

		/* unaproved */
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384, // ECDH-RSA-AES256-SHA384
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
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CipherSuites:             DefaultCipher,
	}

	cert, err := tls.LoadX509KeyPair(tlsCfg.Cert, tlsCfg.Key)
	if err != nil {
		s.log.Fatalf("cannot load cert [%s] and key [%s]: %v", tlsCfg.Cert, tlsCfg.Key, err)
	}

	worker.TLSConfig.Certificates[0] = cert
}
