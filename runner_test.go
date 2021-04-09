package webfmwk

import (
	"encoding/json"
	"net/http"
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
		"\nname: \"oops\"\naddr: \"Testing\"\ntls:\n\tcert:\t\"some/cert\"\n\tkey:\t\"some/key\"\n\tca:\t\"\",\n\tinsecure:\ttrue\n",
		addr.String())
}

func TestRunner(t *testing.T) {
	var s = InitServer(CheckIsUp(), DisableKeepAlive())

	t.Cleanup(func() { stopServer(t, s) })

	s.GET("/test", func(c Context) error {
		return c.JSONOk(json.RawMessage(`{"value":"test"}`))
	})

	go s.Run(Address{Addr: ":6661"}, Address{Addr: ":6662"})

	<-s.isReady

	requestAndTestAPI(t, "http://127.0.0.1:6661/test",
		func(t *testing.T, resp *http.Response) {
			assertStatusCode(t, http.StatusOK, resp)
		})

	requestAndTestAPI(t, "http://127.0.0.1:6662/test",
		func(t *testing.T, resp *http.Response) {
			assertStatusCode(t, http.StatusOK, resp)
		})
}
