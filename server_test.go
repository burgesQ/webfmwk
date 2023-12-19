package webfmwk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDumpRoutes(t *testing.T) {
	s, e := InitServer(
		SetPrefix("/api"),
		CheckIsUp())
	require.Nil(t, e)

	s.GET("/get", func(c Context) error {
		return nil
	})
	s.POST("/post", func(c Context) error {
		return nil
	})
	s.PUT("/put", func(c Context) error {
		return nil
	})
	s.PATCH("/patch", func(c Context) error {
		return nil
	})
	s.DELETE("/delete", func(c Context) error {
		return nil
	})

	all := s.DumpRoutes()

	t.Log(all)

	expected := map[string][]string{
		"DELETE": {"/api/delete"},
		"GET":    {"/api/ping", "/api/get"},
		"PATCH":  {"/api/patch"},
		"POST":   {"/api/post"},
		"PUT":    {"/api/put"},
	}

	require.Equal(t, expected, all)

	// options handled by fasthttp
	// require.Contains(t, all, "OPTIONS")
}

func TestGetLauncher(t *testing.T) {
	s, e := InitServer(CheckIsUp())

	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })
	if s.GetGroup() == nil {
		t.Errorf("Worker group wrongly created: %v.", s.wg)
	}
}

func TestGetContext(t *testing.T) {
	s, e := InitServer(CheckIsUp())

	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

	if s.GetContext() == nil {
		t.Errorf("Context wrongly created : %v.", s.ctx)
	}
}

func TestAddHandlers(t *testing.T) {
	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

	s.addHandlers(func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c Context) error {
			return nil
		})
	})

	require.True(t, len(s.meta.handlers) == 1, "handler wrongly saved")
}

// // // TODO: TestStartTLS(t *testing.T)
// // // TODO: TestStart
// // // TODO: TestShutDown
// // // TODO: TestWaitAndStop
// // // TODO: TestExitHandler
