package webfmwk

import (
	"encoding/json"
)

type (
	// Address hold the api server address
	Address struct {
		Addr string     `json:"addr"`
		TLS  *TLSConfig `json:"tls,omitempty" mapstructure:"tls,omitempty"`
		Name string     `json:"name"`
	}
)

func (a Address) String() string {
	b, e := json.MarshalIndent(a, "	", "	")
	if e != nil {
		return e.Error()
	}
	return string(b)
}

// Run allow to launch multiple server from a single call.
// It take an vaarg Address param argument. WaitAndStop is called via defer.
func (s *Server) Run(addrs ...Address) {
	defer s.WaitAndStop()

	for i := range addrs {
		addr := addrs[i]
		if addr.TLS != nil {
			s.GetLogger().Infof("listening on https://%q", addr.Addr)
			s.StartTLS(addr.Addr, addr.TLS)
		} else {
			s.GetLogger().Infof("listening on http://%q", addr.Addr)
			s.Start(addr.Addr)
		}
	}
}
