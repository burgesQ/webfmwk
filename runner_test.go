package webfmwk

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: test tls

func TestAddress(t *testing.T) {
	addr := new(Address)

	asserter := assert.New(t)

	asserter.Implements((*IAddress)(nil), addr)
	asserter.False(addr.IsOk())

	asserter.Equal("name: \"\"\naddr: \"\"", addr.String())

	addr = &Address{
		Addr: "Testing",
		Name: "oops",
		TLS: &TLSConfig{
			Cert:     "some/cert",
			Key:      "some/key",
			Insecure: true,
		}}

	asserter.Equal("Testing", addr.GetAddr())
	asserter.Equal("oops", addr.GetName())
	asserter.True(addr.IsOk())

	asserter.Equal(
		"name: \"oops\"\naddr: \"Testing\"\ntls: cert:\t\"some/cert\"\nkey:\t\"some/key\"\ninsecure:\ttrue\n",
		addr.String())
}

func TestRunner(t *testing.T) {

	var s = InitServer(CheckIsUp(), DisableKeepAlive())

	t.Log("init server...")
	defer stopServer(t, s)

	s.GET("/test", func(c Context) error {
		return c.JSONOk(json.RawMessage(`{"value":"test"}`))
	})

	go s.Run(Address{Addr: ":6661"}, Address{Addr: ":6662"})

	<-s.isReady
	t.Log("server inited")

	requestAndTestAPI(t, "http://127.0.0.1:6661/test",
		func(t *testing.T, resp *http.Response) {
			assertStatusCode(t, http.StatusOK, resp)
		})

	requestAndTestAPI(t, "http://127.0.0.1:6662/test",
		func(t *testing.T, resp *http.Response) {
			assertStatusCode(t, http.StatusOK, resp)
		})
}
