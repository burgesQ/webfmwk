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

	// Returned in case of invalid ca cert path.
	ErrParseUserCA = errors.New("failed to parse root certificate")
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

	listener, err := LoadTLSListener(addr, tlsCfg)
	if err != nil {
		s.GetLogger().Fatalf("loading tls: %s", err.Error())
	}

	s.launcher.Start("https server "+addr, func() error {
		return s.internalInit(addr).Serve(listener)
	})
}

// LoadTLSListener return a tls listner ready for mTLS.
func LoadTLSListener(addr string, tlsCfg ITLSConfig) (net.Listener, error) {
	var cfg, e = GetTLSCfg(tlsCfg)
	if e != nil {
		return nil, e
	}

	listner, e := net.Listen("tcp4", addr)
	if e != nil {
		return nil, e
	}

	return tls.NewListener(listner, cfg), nil
}

// GetTLSCfg return a tls config ready for mTLS.
// thx to https://dev.to/living_syn/validating-client-certificate-sans-in-go-i5p
func GetTLSCfg(tlsCfg ITLSConfig) (*tls.Config, error) {
	var cert, err = tls.LoadX509KeyPair(tlsCfg.GetCert(), tlsCfg.GetKey())
	if err != nil {
		return nil, fmt.Errorf("cannot load cert [%s] and key [%s]: %w",
			tlsCfg.GetCert(), tlsCfg.GetKey(), err)
	}

	/* #nosec */
	cfg := getBaseTLSCfg(&cert)

	if tlsCfg.GetInsecure() {
		cfg.ClientAuth = tls.NoClientCert

		return cfg, nil
	}

	cfg.ClientAuth = tls.RequireAndVerifyClientCert
	if e := loadCA(tlsCfg.GetCa(), cfg); e != nil {
		return cfg, e
	}

	cfg.GetConfigForClient = wrapGetConfigForClient(&cert, cfg.ClientCAs)

	return cfg, nil
}

func loadCA(caPath string, cfg *tls.Config) error {
	if caPath == "" {
		return nil
	}

	pool := x509.NewCertPool()

	if caCertPEM, e := ioutil.ReadFile(caPath); e != nil {
		return fmt.Errorf("cannot load ca cert %q in pool: %w", caPath, e)
	} else if !pool.AppendCertsFromPEM(caCertPEM) {
		return ErrParseUserCA
	}

	cfg.ClientCAs = pool

	return nil
}

func getBaseTLSCfg(cert *tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates:             []tls.Certificate{*cert},
		PreferServerCipherSuites: true,
		CurvePreferences:         DefaultCurve,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13, // tls.VersionTLS12 ?
		CipherSuites:             DefaultCipher,
	}
}

func wrapVerifyPerrCertificate(caCert *x509.CertPool, remoteAddr string) func(
	[][]byte, [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// from src/crypto/tls/handshake_server.go:680 (go 1.11) + DNSName check
		var opts = x509.VerifyOptions{
			Roots:         caCert,
			CurrentTime:   time.Now(),
			Intermediates: x509.NewCertPool(),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			DNSName:       strings.Split(remoteAddr, ":")[0],
		}

		_, err := verifiedChains[0][0].Verify(opts)

		return err
	}
}

func wrapGetConfigForClient(cert *tls.Certificate, caCert *x509.CertPool) func(
	*tls.ClientHelloInfo) (*tls.Config, error) {
	return func(hi *tls.ClientHelloInfo) (*tls.Config, error) {
		var cfg = getBaseTLSCfg(cert)

		cfg.ClientAuth, cfg.ClientCAs = tls.RequireAndVerifyClientCert, caCert
		cfg.VerifyPeerCertificate = wrapVerifyPerrCertificate(caCert,
			hi.Conn.RemoteAddr().String())

		return cfg, nil
	}
}
