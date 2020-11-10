package webfmwk

import (
	"net/http"
	"time"
)

func (s *Server) getResp(uri string) (r *http.Response, e error) {
	req, err := http.NewRequestWithContext(s.ctx, "GET", uri, nil)
	if err != nil {
		return r, err
	}

	client := http.DefaultClient
	r, e = client.Do(req)

	return
}

// pollPingEndpoint try to reach the /ping endpoint of the server
// to then infrome that the server is up via the isReady channel
func (s *Server) pollPingEndpoint(addr string) {
	if !s.meta.checkIsUp {
		return
	}

	var (
		uri   = concatAddr(addr, s.meta.prefix)
		delay = time.Millisecond * 0
	)

	defer func() {
		s.isReady <- true
	}()

	for s.ctx.Err() == nil {
		select {
		case <-time.After(delay):
			delay = time.Millisecond * 10
			/* #nosec  */
			if resp, e := s.getResp(uri); e != nil {
				s.log.Infof("server not up (%q) ... %s", uri, e.Error())
				continue
			} else if e = resp.Body.Close(); e != nil || resp.StatusCode != http.StatusOK {
				s.log.Infof("unexpected status code, %s : %v", resp.StatusCode, e)
				continue
			}

			s.log.Infof("server is up")
			return

		case <-s.ctx.Done():
			return
		}
	}
}
