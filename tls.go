package webfmwk

import (
	"crypto/tls"
	"fmt"
	"net"
)

type (
	// ITLSConfig is used to interface the TLS implemtation.
	ITLSConfig interface {
		fmt.Stringer

		// GetCert return the full path to the server certificate file.
		GetCert() string

		// GetKey return the full path to the server key file.
		GetKey() string

		// GetInsecure return true if the TLS Certificate shouldn't be checked.
		GetInsecure() bool

		// IsEmpty return true if the config is empty.
		Empty() bool
	}

	// TLSConfig contain the tls config passed by the config file.
	// It implement TLSConfig
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

// GetCert implemte TLSConfig
func (config TLSConfig) GetCert() string {
	return config.Cert
}

// GetKey implemte TLSConfig
func (config TLSConfig) GetKey() string {
	return config.Key
}

// GetInsecure implemte TLSConfig
func (config TLSConfig) GetInsecure() bool {
	return config.Insecure
}

// Empty implemte TLSConfig
func (config TLSConfig) Empty() bool {
	return config.Cert == "" && config.Key == "" && !config.Insecure
}

// String implement Stringer interface
func (config TLSConfig) String() string {
	if config.Empty() {
		return ""
	}

	return fmt.Sprintf("cert:\t%q\nkey:\t%q\ninsecure:\t%t\n",
		config.Cert, config.Key, config.Insecure)
}

// StartTLS expose an server to an HTTPS address.
func (s *Server) StartTLS(addr string, tlsStuffs ITLSConfig) {
	s.internalHandler()
	s.launcher.Start("https server "+addr, func() error {
		return s.internalInit(addr).Serve(s.loadTLSListener(addr, tlsStuffs))
	})
}

func (s *Server) getTLSCfg(tlsCfg ITLSConfig) *tls.Config {
	cert, err := tls.LoadX509KeyPair(tlsCfg.GetCert(), tlsCfg.GetKey())
	if err != nil {
		s.log.Fatalf("cannot load cert [%s] and key [%s]: %s",
			tlsCfg.GetCert(), tlsCfg.GetKey(), err.Error())
	}

	/* #nosec */
	return &tls.Config{
		InsecureSkipVerify:       tlsCfg.GetInsecure(),
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CipherSuites:             DefaultCipher,
	}
}

func (s *Server) loadTLSListener(addr string, tlsCfg ITLSConfig) net.Listener {
	cfg := s.getTLSCfg(tlsCfg)

	tmpLn, err := net.Listen("tcp4", addr)
	if err != nil {
		s.log.Fatalf("cannot listen on %q: %s", addr, err.Error())
	}

	return tls.NewListener(tmpLn, cfg)
}
