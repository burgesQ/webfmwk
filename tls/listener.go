package tls

import (
	"crypto/tls"
	"fmt"
	"net"
)

// LoadTLSListener return a tls listner ready for mTLS and/or http2.
func LoadListner(addr string, cfg *tls.Config) (net.Listener, error) {
	listner, e := net.Listen("tcp4", addr)
	if e != nil {
		return nil, fmt.Errorf("creating tls listner: %w", e)
	}

	return tls.NewListener(listner, cfg), nil
}
