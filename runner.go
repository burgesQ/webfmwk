package webfmwk

import (
	"encoding/json"
)

type (
	IAddress interface {
		GetAddr() string
		GetTLS() *TLSConfig
		GetName() string
		IsOk() bool
	}

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

func (a Address) IsOk() bool {
	return a.Addr != ""
}

func (a Address) GetAddr() string {
	return a.Addr
}

func (a Address) GetTLS() *TLSConfig {
	return a.TLS
}

func (a Address) GetName() string {
	return a.Name
}

// Run allow to launch multiple server from a single call.
// It take an vaarg Address param argument. WaitAndStop is called via defer.
func (s *Server) Run(addrs ...Address) {
	defer s.WaitAndStop()

	for i := range addrs {
		addr := addrs[i]
		if addr.IsOk() {
			if tls := addr.GetTLS(); tls != nil {
				s.GetLogger().Infof("starting %s on https://%s", addr.GetName(), addr.GetAddr())
				s.StartTLS(addr.GetAddr(), tls)
			} else {
				s.GetLogger().Infof("starting %s on http://%s", addr.GetName(), addr.GetAddr())
				s.Start(addr.GetAddr())
			}
		} else {
			s.GetLogger().Errorf("invalid address format : %s", addr)
		}
	}
}
