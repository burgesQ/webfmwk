package webfmwk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type (
	// ITLSConfig is used to interface the TLS implemtation.
	ITLSConfig interface {
		fmt.Stringer

		// GetCert return the full path to the server certificate file.
		GetCert() string

		// GetKey return the full path to the server key file.
		GetKey() string

		// GetCa return the full path to the server ca cert file.
		GetCa() string

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
		Ca       string `json:"ca" mapstructur:"ca"`
		Insecure bool   `json:"insecure" mapstructur:"insecure"`
	}
)

var (
	// DefaultCurve TLS curve supported
	DefaultCurve = []tls.CurveID{
		tls.CurveP256,
		tls.X25519,
		// tls.CurveP384
		// tls.CurveP521

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

// GetCert implemte TLSConfig.
func (config TLSConfig) GetCert() string {
	return config.Cert
}

// GetKey implemte TLSConfig.
func (config TLSConfig) GetKey() string {
	return config.Key
}

// GetKey implemte TLSConfig.
func (config TLSConfig) GetCa() string {
	return config.Ca
}

// GetInsecure implemte TLSConfig.
func (config TLSConfig) GetInsecure() bool {
	return config.Insecure
}

// Empty implemte TLSConfig.
func (config TLSConfig) Empty() bool {
	return config.Cert == "" && config.Key == ""
}

// String implement Stringer interface.
func (config TLSConfig) String() string {
	if config.Empty() {
		return ""
	}

	return fmt.Sprintf("\t ~!~ cert:\t%q\n\t ~!~ key:\t%q\n\t ~!~ ca:\t%q,\n\t ~!~ insecure:\t%t\n",
		config.Cert, config.Key, config.Ca, config.Insecure)
}

// StartTLS expose an server to an HTTPS address..
func (s *Server) StartTLS(addr string, tlsCfg ITLSConfig) {
	s.internalHandler()

	listener := s.loadTLSListener(addr, tlsCfg)

	s.launcher.Start("https server "+addr, func() error {
		return s.internalInit(addr).Serve(listener)
	})
}

func (s *Server) getTLSCfgCA(tlsCfg ITLSConfig) (*tls.Config, error) {
	var lg = GetLogger()

	cert, err := tls.LoadX509KeyPair(tlsCfg.GetCert(), tlsCfg.GetKey())
	if err != nil {
		lg.Fatalf("cannot load cert [%s] and key [%s]: %s",
			tlsCfg.GetCert(), tlsCfg.GetKey(), err.Error())
	}

	/* #nosec */
	cfg := &tls.Config{
		// ServerName:               tlsCfg.GetName(),
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13, // tls.VersionTLS12 ?
		CipherSuites:             DefaultCipher,
	}

	if tlsCfg.GetInsecure() {
		cfg.ClientAuth = tls.ClientAuthType(tls.NoClientCert)
		lg.Debugf("no cert req")
	} else {
		lg.Debugf("req and verify")
		cfg.ClientAuth = tls.ClientAuthType(tls.RequireAndVerifyClientCert)

		if caPath := tlsCfg.GetCa(); caPath != "" {
			lg.Debugf("no path")

			pool := x509.NewCertPool()
			if caCertPEM, e := ioutil.ReadFile(caPath); e != nil {
				return cfg, fmt.Errorf("cannot load ca cert pool | %w", caPath, e)
			} else if !pool.AppendCertsFromPEM(caCertPEM) {
				return cfg, errors.New("failed to parse root certificate")
			}

			cfg.ClientCAs = pool

			cfg.BuildNameToCertificate()
		}
	}

	cfg.GetConfigForClient = func(hi *tls.ClientHelloInfo) (*tls.Config, error) {
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    cfg.ClientCAs,
			VerifyPeerCertificate: func(helloInfo *tls.ClientHelloInfo) func([][]byte, [][]*x509.Certificate) error {
				return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					//copied from the default options in src/crypto/tls/handshake_server.go, 680 (go 1.11)
					//but added DNSName
					opts := x509.VerifyOptions{
						Roots:         cfg.ClientCAs,
						CurrentTime:   time.Now(),
						Intermediates: x509.NewCertPool(),
						KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
						DNSName:       strings.Split(helloInfo.Conn.RemoteAddr().String(), ":")[0],
					}
					_, err := verifiedChains[0][0].Verify(opts)
					return err

					// return nil
				}
			}(hi),
		}, nil
	}

	return cfg, nil
}

func (s *Server) loadTLSListener(addr string, tlsCfg ITLSConfig) net.Listener {
	var (
		cfg, e = s.getTLSCfgCA(tlsCfg)
	)

	if e != nil {
		s.GetLogger().Fatalf("%s", e.Error())
	}

	listner, err := net.Listen("tcp4", addr)
	if err != nil {
		s.GetLogger().Fatalf("cannot listen on %q: %s", addr, err.Error())
	}

	return tls.NewListener(listner, cfg)
}

func (s *Server) getTLSCfg(tlsCfg ITLSConfig) *tls.Config {
	cert, err := tls.LoadX509KeyPair(tlsCfg.GetCert(), tlsCfg.GetKey())
	if err != nil {
		s.GetLogger().Fatalf("cannot load cert [%s] and key [%s]: %s",
			tlsCfg.GetCert(), tlsCfg.GetKey(), err.Error())
	}

	/* #nosec */
	return &tls.Config{
		// ServerName:               tlsCfg.GetName(),
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13, // tls.VersionTLS12 ?
		CipherSuites:             DefaultCipher,
	}
}

// mtls thx to https://dev.to/living_syn/validating-client-certificate-sans-in-go-i5p

// register ca cert pool and toggle cert requirement
func loadCa(cfg *tls.Config, insecure bool, caPath string) error {
	lg := GetLogger()

	if insecure {
		cfg.ClientAuth = tls.ClientAuthType(tls.NoClientCert)
		lg.Debugf("no cert req")

		return nil
	}

	lg.Debugf("req and verify")
	cfg.ClientAuth = tls.ClientAuthType(tls.RequireAndVerifyClientCert)

	if caPath == "" {
		lg.Debugf("no path")
		return nil
	}

	pool := x509.NewCertPool()
	if caCertPEM, e := ioutil.ReadFile(caPath); e != nil {
		return fmt.Errorf("cannot load ca cert pool | %s", caPath, e.Error())
	} else if !pool.AppendCertsFromPEM(caCertPEM) {
		return fmt.Errorf("failed to parse root certificate")
	}

	cfg.ClientCAs = pool

	return nil
}
