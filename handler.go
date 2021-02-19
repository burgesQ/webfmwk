package webfmwk

import (
	"errors"
	"net/http"
)

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	s.log.Infof("[!] 404 reached for [%s] %s %s", GetIPFromRequest(r), r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if _, e := w.Write([]byte(`{"status":404,"message":"not found"}`)); e != nil {
		s.log.Errorf("[!] cannot write 404 ! %s", e.Error())
	}
}

func (s *Server) handleNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.log.Infof("[!] 405 reached for [%s] %s %s", GetIPFromRequest(r), r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	if _, e := w.Write([]byte(`{"status":405,"message":"method not allowed"}`)); e != nil {
		s.log.Errorf("cannot write 405 ! %s", e.Error())
	}
}

func (s *Server) handleError(ctx Context, e error) {
	var eh ErrorHandled
	if errors.As(e, &eh) {
		_ = ctx.JSON(eh.GetOPCode(), eh.GetContent())
		return
	}

	_ = ctx.JSONInternalError(NewErrorFromError(e))
}
