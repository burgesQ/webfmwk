package webfmwk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: test tls

func TestAddress(t *testing.T) {
	addr := new(Address)

	requirer := require.New(t)

	requirer.Implements((*IAddress)(nil), addr)
	requirer.False(addr.IsOk())

	requirer.Equal("name: \"\"\naddr: \"\"", addr.String())

	addr = &Address{
		Addr: "Testing",
		Name: "oops",
		TLS: &TLSConfig{
			Cert:     "some/cert",
			Key:      "some/key",
			Insecure: true,
		}}

	requirer.Equal("Testing", addr.GetAddr())
	requirer.Equal("oops", addr.GetName())
	requirer.True(addr.IsOk())

	requirer.Equal(
		"\nname: \"oops\"\naddr: \"Testing\"\ntls:\n\tcert:\t\"some/cert\"\n"+
			"\tkey:\t\"some/key\"\n\tca:\t\"\",\n\tinsecure:\ttrue\n",
		addr.String())

	requirer.Equal(
		"number of address(es): 2\n\nname: \"oops\"\naddr: \"Testing\"\ntls:\n"+
			"\tcert:\t\"some/cert\"\n\tkey:\t\"some/key\"\n\tca:\t\"\",\n"+
			"\tinsecure:\ttrue\n\nname: \"smth\"\naddr: \"Testing_2\"",
		Addresses{
			*addr,
			Address{
				Addr: "Testing_2",
				Name: "smth",
			},
		}.String())
}
