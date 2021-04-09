package webfmwk

import (
	"testing"

	"github.com/burgesQ/webfmwk/v5/log"
	"github.com/stretchr/testify/assert"
)

func TestDumpRoutes(t *testing.T) {

	s := InitServer(
		SetPrefix("/api"),
		CheckIsUp())

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
		"DELETE": []string{"/api/delete"},
		"GET":    []string{"/api/ping", "/api/get"},
		"PATCH":  []string{"/api/patch"},
		"POST":   []string{"/api/post"},
		"PUT":    []string{"/api/put"},
	}

	assert.Equal(t, expected, all)

	// options handled by fasthttp
	// assert.Contains(t, all, "OPTIONS")

}

type customLoggerT struct{}

func (l customLoggerT) Printf(format string, args ...interface{}) {}
func (l customLoggerT) Debugf(format string, v ...interface{})    {}
func (l customLoggerT) Infof(format string, v ...interface{})     {}
func (l customLoggerT) Warnf(format string, v ...interface{})     {}
func (l customLoggerT) Errorf(format string, v ...interface{})    {}
func (l customLoggerT) Fatalf(format string, v ...interface{})    {}

func TestRegisterLogger(t *testing.T) {
	var (
		lg = new(customLoggerT)
		s  = InitServer(WithLogger(lg))
	)

	assert.Implements(t, (*log.Log)(nil), lg)
	assert.Equal(t, lg, s.GetLogger())
}

func TestGetLauncher(t *testing.T) {
	s := InitServer(CheckIsUp())

	t.Cleanup(func() { stopServer(t, s) })
	if s.GetLauncher() == nil {
		t.Errorf("Launcher wrongly created : %v.", s.launcher)
	}
}

func TestGetContext(t *testing.T) {
	s := InitServer(CheckIsUp())

	t.Cleanup(func() { stopServer(t, s) })

	if s.GetContext() == nil {
		t.Errorf("Context wrongly created : %v.", s.ctx)
	}
}

func TestAddHandlers(t *testing.T) {
	s := InitServer(CheckIsUp())
	t.Cleanup(func() { stopServer(t, s) })

	s.addHandlers(func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c Context) error {
			return nil
		})
	})

	assert.True(t, len(s.meta.handlers) == 1, "handler wrongly saved")
}

// // // TODO: TestStartTLS(t *testing.T)
// // // TODO: TestStart
// // // TODO: TestShutDown
// // // TODO: TestWaitAndStop
// // // TODO: TestExitHandler
