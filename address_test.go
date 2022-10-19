package webfmwk

import (
	"testing"

	"github.com/burgesQ/webfmwk/v5/tls"
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
		TLS: &tls.Config{
			Cert:     "some/cert",
			Key:      "some/key",
			Insecure: true,
		}}

	requirer.Equal("Testing", addr.GetAddr())
	requirer.Equal("oops", addr.GetName())
	requirer.True(addr.IsOk())

	requirer.Equal(
		"\n\t -!- name: \"oops\"\n\t -!- addr: \"Testing\"\n"+
			"\t -!- tls:\n\t ~!~ cert:\t\"some/cert\"\n"+
			"\t ~!~ key:\t\"some/key\"\n\t ~!~ ca:\t\"\",\n"+
			"\t ~!~ insecure:\ttrue\n\t ~!~ level:\tnever\n",
		addr.String())

	requirer.Equal(

		"\t --- number of address(es): 2\n\n\t -!- name: \"oops\"\n"+
			"\t -!- addr: \"Testing\"\n\t -!- tls:\n"+
			"\t ~!~ cert:\t\"some/cert\"\n\t ~!~ key:\t\"some/key\"\n"+
			"\t ~!~ ca:\t\"\",\n\t ~!~ insecure:\ttrue\n"+
			"\t ~!~ level:\tnever\n"+
			"\nname: \"smth\"\naddr: \"Testing_2\"\n\t --- end address",
		Addresses{
			*addr,
			Address{
				Addr: "Testing_2",
				Name: "smth",
			},
		}.String())
}

func TestAddressSameAs(t *testing.T) {
	addr := &Address{
		Addr: "uno",
		Name: "deuzio",
		TLS: &tls.Config{
			Cert:     "some/cert",
			Key:      "some/key",
			Insecure: true,
		},
	}
	requirer := require.New(t)

	requirer.False(addr.SameAs(&Address{}))
	requirer.False(addr.SameAs(&Address{
		Addr: "uno",
		Name: "deuzio",
	}))

	requirer.False(addr.SameAs(&Address{
		Addr: "uno",
		Name: "deuzio",
		TLS: &tls.Config{
			Cert:     "some/cert",
			Insecure: true,
		}}))

	requirer.True(addr.SameAs(&Address{
		Addr: "uno",
		Name: "deuzio",
		TLS: &tls.Config{
			Cert:     "some/cert",
			Key:      "some/key",
			Insecure: true,
		}}))
}
