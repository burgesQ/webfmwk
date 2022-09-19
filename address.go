package webfmwk

import (
	"fmt"
	"strconv"
)

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

	Addresses []Address
)

func (a Addresses) String() (ret string) {
	ret = "\t --- number of address(es): " + strconv.Itoa(len(a))

	for i := range a {
		ret += "\n" + a[i].String()
	}

	ret += "\n\t --- end address"

	return
}

// String implement the fmt.Stringer interface
func (a Address) String() string {
	if a.TLS != nil && !a.TLS.Empty() {
		return fmt.Sprintf("\n\t -!- name: %q\n\t -!- addr: %q\n\t -!- tls:\n%s", a.Name, a.Addr, a.TLS.String())
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
