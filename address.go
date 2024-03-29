package webfmwk

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/burgesQ/webfmwk/v6/tls"
)

const _unixSocketPrefix = `unix://`

type (
	// IAddress interface hold an api server listing configuration
	IAddress interface {
		fmt.Stringer

		// GetAddr return the listing address.
		GetAddr() string

		// GetTLS return a pointer to an TLSConfig if present, nil otherwise.
		GetTLS() tls.IConfig

		// GetName return the name of the address, for debug purpose.
		GetName() string

		// IsOk validate that the Address structure have at least the address field populated.
		IsOk() bool

		// IsUnixPath return true if the address is a valid unix socket path.
		IsUnixPath() bool

		// SameAs return true if both config are identique.
		SameAs(in IAddress) bool
	}

	// Address implement the IAddress interface
	Address struct {
		// TLS implement IAddress, tlsConfig  implement the TLSConfig interface.
		TLS  *tls.Config `json:"tls,omitempty" mapstructure:"tls,omitempty"`
		Addr string      `json:"addr"`
		Name string      `json:"name"`
	}

	Addresses []Address
)

func (a Address) SameAs(in IAddress) bool {
	var (
		tlsOk = false
		itls  = in.GetTLS()
	)

	if a.TLS == nil && itls == nil {
		tlsOk = true
	} else if a.TLS != nil && itls != nil {
		tlsOk = a.TLS.SameAs(itls)
	}

	return a.Addr == in.GetAddr() && a.Name == in.GetName() && tlsOk
}

// SameAs return true if all addresses in the in param match
// One from the struct.
func (a Addresses) SameAs(in Addresses) bool {
	if len(a) != len(in) { // not same nb of address
		return false
	}

	for i := range a {
		iok := false

		for j := range in {
			if in[j].SameAs(a[i]) {
				iok = true

				break
			}
		}

		if !iok {
			return false
		}
	}

	return true
}

func (a Addresses) String() (ret string) {
	ret = "\t --- number of address(es): " + strconv.Itoa(len(a))

	for i := range a {
		ret += "\n" + a[i].String()
	}

	ret += "\n\t --- end address"

	return
}

func (a Addresses) AsAttrs() []any {
	r := make([]any, len(a))

	for i := range a {
		r[i] = slog.Group(fmt.Sprintf("address %d", i), a[i].AsAttrs()...)
	}

	return r
}

func Tern[T any](cond bool, t, f func() T) T {
	if cond {
		return t()
	}

	return f()
}

func (a Address) AsAttrs() []any {
	return []any{
		slog.String("name", a.Name),
		slog.String("address", a.Addr),
		slog.Any("tls", a.TLS),
	}
}

// String implement the fmt.Stringer interface
func (a Address) String() string {
	if a.TLS != nil && !a.TLS.Empty() {
		return fmt.Sprintf("\n\t -!- name: %q\n\t -!- addr: %q\n\t -!- tls:\n%s", a.Name, a.Addr, a.TLS.String())
	}

	return fmt.Sprintf("name: %q\naddr: %q", a.Name, a.Addr)
}

// IsOk implement the IAddress interface
func (a Address) IsUnixPath() bool {
	return strings.HasPrefix(a.Addr, _unixSocketPrefix)
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
func (a Address) GetTLS() tls.IConfig {
	if a.TLS == nil {
		return nil
	}

	return *a.TLS
}

// GetName implement the IAddress interface
func (a Address) GetName() string {
	return a.Name
}
