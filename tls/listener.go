package tls

import (
	"crypto/tls"
	"net"
)

// LoadTLSListener return a tls listner ready for mTLS.
func LoadListener(addr string, icfg IConfig) (net.Listener, error) {
	cfg, e := GetTLSCfg(icfg)
	if e != nil {
		return nil, e
	}

	listner, e := net.Listen("tcp4", addr)
	if e != nil {
		return nil, e
	}

	return tls.NewListener(listner, cfg), nil
}
