package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
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

// GetTLSCfg return a tls config ready for mTLS.
// thx to https://dev.to/living_syn/validating-client-certificate-sans-in-go-i5p
func GetTLSCfg(icfg IConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(icfg.GetCert(), icfg.GetKey())
	if err != nil {
		return nil, fmt.Errorf("cannot load cert [%s] and key [%s]: %w",
			icfg.GetCert(), icfg.GetKey(), err)
	}

	/* #nosec */
	cfg := getBaseTLSCfg(&cert)

	if icfg.GetInsecure() {
		cfg.ClientAuth = tls.NoClientCert

		return cfg, nil
	}

	lvl := icfg.GetLevel()
	if lvl == NoClientCert {
		lvl = RequestClientCert
	}

	cfg.ClientAuth = lvl.STD()
	if e := loadCA(icfg.GetCa(), cfg); e != nil {
		return cfg, e
	}

	cfg.GetConfigForClient = wrapGetConfigForClient(&cert, cfg.ClientCAs, lvl)

	return cfg, nil
}

func loadCA(caPath string, cfg *tls.Config) error {
	if caPath == "" {
		return nil
	}

	pool := x509.NewCertPool()

	if caCertPEM, e := os.ReadFile(caPath); e != nil {
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

func wrapGetConfigForClient(cert *tls.Certificate, caCert *x509.CertPool,
	level Level,
) func(*tls.ClientHelloInfo) (*tls.Config, error) {
	return func(hi *tls.ClientHelloInfo) (*tls.Config, error) {
		cfg := getBaseTLSCfg(cert)

		cfg.ClientAuth, cfg.ClientCAs = level.STD(), caCert
		if level == RequireAndVerifyClientCertAndSAN {
			cfg.VerifyPeerCertificate = wrapVerifyPerrCertificate(caCert,
				hi.Conn.RemoteAddr().String())
		}

		return cfg, nil
	}
}

func wrapVerifyPerrCertificate(caCert *x509.CertPool, remoteAddr string) func(
	[][]byte, [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// from src/crypto/tls/handshake_server.go:680 (go 1.11) + DNSName check
		opts := x509.VerifyOptions{
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
