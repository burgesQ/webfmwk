package webfmwk

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	// ITLSConfig is used to interface the TLS implemtation.
	ITLSConfig interface {
		fmt.Stringer

		// GetCert return the full path to the server certificate file
		GetCert() string
		// GetKey return the full path to the server key file
		GetKey() string
		// GetInsecure return true if the TLS Certificate shouldn't be checked
		GetInsecure() bool
		// IsEmpty return true if the config is empty
		Empty() bool
	}

	// TLSConfig contain the tls config passed by the config file.
	// It implement ITLSConfig
	TLSConfig struct {
		Cert     string `json:"cert" mapstructur:"cert"`
		Key      string `json:"key" mapstructur:"key"`
		Insecure bool   `json:"insecure" mapstructur:"insecure"`
	}
)

var (
	// DefaultCurve TLS curve supported
	DefaultCurve = []tls.CurveID{
		tls.CurveP256,
		tls.X25519,
	}
	// DefaultCipher accepted
	DefaultCipher = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,   // HTTP/2-required AES_128_GCM_SHA256 cipher
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   // ECDHE-RSA-AES256-GCM-SHA384
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   // ECDH-RSA-AES256-GCM-SHA384
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, // ECDHE-ECDSA-AES256-GCM-SHA384
		tls.TLS_AES_128_GCM_SHA256,                  // 1.3 tls cipher
		tls.TLS_AES_256_GCM_SHA384,                  // 1.3 tls cipher
		tls.TLS_CHACHA20_POLY1305_SHA256,            // 1.3 tls cipher

		/* unaproved ? */
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384, // ECDH-RSA-AES256-SHA384
	}
)

// GetCert implemte ITLSConfig
func (config TLSConfig) GetCert() string {
	return config.Cert
}

// GetKey implemte ITLSConfig
func (config TLSConfig) GetKey() string {
	return config.Key
}

// GetInsecure implemte ITLSConfig
func (config TLSConfig) GetInsecure() bool {
	return config.Insecure
}

// GetInsecure implemte ITLSConfig
func (config TLSConfig) Empty() bool {
	return len(config.Cert) == 0 && len(config.Key) == 0 && !config.Insecure
}

// String implement Stringer interface
func (config TLSConfig) String() string {
	b, e := json.MarshalIndent(config, " ", "\t")
	if e != nil {
		return "error"
	}

	return string(b)
}

// StartTLS expose an server to an HTTPS address
func (s *Server) StartTLS(addr string, tlsStuffs ITLSConfig) {
	s.internalHandler()
	s.launcher.Start("https server "+addr, func() error {
		go s.pollPingEndpoint(addr)
		return s.internalInit(addr, tlsStuffs).ListenAndServeTLS(tlsStuffs.GetCert(), tlsStuffs.GetKey())
	})
}

func (s *Server) loadTLS(worker *http.Server, tlsCfg ITLSConfig) {
	cert, err := tls.LoadX509KeyPair(tlsCfg.GetCert(), tlsCfg.GetKey())
	if err != nil {
		s.log.Fatalf("cannot load cert [%s] and key [%s]: %v", tlsCfg.GetCert(), tlsCfg.GetKey(), err)
	}

	/* #nosec */
	worker.TLSConfig = &tls.Config{
		InsecureSkipVerify:       tlsCfg.GetInsecure(),
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CipherSuites:             DefaultCipher,
	}
}
