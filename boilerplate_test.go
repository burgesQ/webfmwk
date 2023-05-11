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

func stopServer(s *Server) {
	ctx := s.GetContext()

	ctx.Done()
	s.Shutdown()
	s.WaitAndStop()
	Shutdown()
}

func wrapperPost(t *testing.T, route, routeReq string,
	content []byte,
	handlerRoute HandlerFunc, handlerTest webtest.HandlerForTest,
) {
	t.Helper()

	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)

	t.Cleanup(func() { stopServer(s) })
	s.POST(route, handlerRoute)
	go s.Start(_testPort)
	<-s.isReady

	webtest.PushAndTestAPI(t, _testAddr+routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute HandlerFunc, handlerTest webtest.HandlerForTest,
) {
	t.Helper()

	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)

	t.Cleanup(func() { stopServer(s) })

	s.GET(route, handlerRoute)
	go s.Start(_testPort)
	<-s.isReady

	webtest.RequestAndTestAPI(t, _testAddr+routeReq, handlerTest)
}
