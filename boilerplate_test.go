package webfmwk

import (
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/stretchr/testify/require"
)

const (
	_testPort = ":6666"
	_testAddr = "http://127.0.0.1" + _testPort
)

//
// helpers methods
//

func wrapperPost(t *testing.T, route, routeReq string,
	content []byte,
	handlerRoute HandlerFunc, handlerTest webtest.HandlerForTest,
) {
	t.Helper()

	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)

	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })
	s.POST(route, handlerRoute)

	go require.Nil(t, s.Start(_testPort))

	<-s.isReady

	webtest.PushAndTestAPI(t, _testAddr+routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute HandlerFunc, handlerTest webtest.HandlerForTest,
) {
	t.Helper()

	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)

	t.Log("server is created")
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })
	s.GET(route, handlerRoute)
	t.Logf("starting server -- %q\n", route)

	go require.Nil(t, s.Start(_testPort))

	<-s.isReady

	t.Logf("server's ready\nrequesting the API on %q\n", _testAddr+routeReq)

	webtest.RequestAndTestAPI(t, _testAddr+routeReq, handlerTest)
}
