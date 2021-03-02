package webfmwk

import (
	"net/http"
)

func GetIPFromRequest(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	return r.RemoteAddr
}

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
