package webfmwk

import "fmt"

type (
	// IAddress interface hold an api server listing configuration
	IAddress interface {
		fmt.Stringer

		// GetAddr return the listing address.
		GetAddr() string

		// GetTLS return a pointer to an TLSConfig if present, nil otherwise.
		GetTLS() ITLSConfig

		// GetName return the name of the address, for debug purpose.
		GetName() string

		// IsOk validate that the Address structure have at least the address field populated.
		IsOk() bool
	}

	// Address implement the IAddress interface
	Address struct {
		Addr string `json:"addr"`
		Name string `json:"name"`

		// TLS implement IAddress, tlsConfig  implement the TLSConfig interface.
		TLS *TLSConfig `json:"tls,omitempty" mapstructure:"tls,omitempty"`
	}
)

// String implement the fmt.Stringer interface
func (a Address) String() string {
	if a.TLS != nil && !a.TLS.Empty() {
		return fmt.Sprintf("name: %q\naddr: %q\ntls: %s", a.Name, a.Addr, a.TLS.String())
	}

	return fmt.Sprintf("name: %q\naddr: %q", a.Name, a.Addr)
}

// IsOk implement the IAddress interface
func (a Address) IsOk() bool {
	return a.Addr != ""
}

// GetAddr implement the IAddress interface
func (a Address) GetAddr() string {
	return a.Addr
}

// GetTLS implement the IAddress interface
func (a Address) GetTLS() ITLSConfig {
	if a.TLS == nil {
		return nil
	}

	return *a.TLS
}

// GetName implement the IAddress interface
func (a Address) GetName() string {
	return a.Name
}

// Run allow to launch multiple server from a single call.
// It take an va arg list of Address as argument.
// The method wait for the server to end via a call to WaitAndStop.
func (s *Server) Run(addrs ...Address) {
	defer s.WaitAndStop()

	for i := range addrs {
		addr := addrs[i]
		if !addr.IsOk() {
			s.GetLogger().Errorf("invalid address format : %s", addr)
			continue
		}

		if tls := addr.GetTLS(); tls != nil && !tls.Empty() {
			s.GetLogger().Infof("starting %s on https://%s", addr.GetName(), addr.GetAddr())
			s.StartTLS(addr.GetAddr(), tls)

			continue
		}

		s.GetLogger().Infof("starting %s on http://%s", addr.GetName(), addr.GetAddr())
		s.Start(addr.GetAddr())
	}
}
