package webfmwk

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/burgesQ/webfmwk/v6/tls"
)

func getResp(ctx context.Context, uri string, cfg ...tls.IConfig) (*http.Response, error) {
	if len(cfg) > 0 {
		uri = "https://" + strings.TrimPrefix(uri, "http://")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient

	return client.Do(req)
}

// pollPingEndpoint try to reach the /ping endpoint of the server
// to then infrome that the server is up via the isReady channel
func (s *Server) pollPingEndpoint(addr string, cfg ...tls.IConfig) {
	var (
		uri      = concatAddr(addr, s.meta.prefix)
		duration = time.Millisecond * 10
	)

	if !s.meta.checkIsUp {
		return
	}

	defer func() {
		s.slog.Info("server is up")
		s.isReady <- true
	}()

	delay := time.NewTimer(time.Millisecond * 0)
	defer delay.Stop()

	for s.ctx.Err() == nil {
		delay.Reset(duration)
		select {
		case <-delay.C:
			delay.Reset(duration)
			/* #nosec  */
			resp, e := getResp(s.ctx, uri, cfg...)
			if e != nil {
				str := e.Error()

				if strings.HasSuffix(str, "unknown authority") ||
					strings.HasSuffix(str, "any IP SANs") ||
					errors.Is(e, io.EOF) {
					return
				}

				s.slog.Info("server not up",
					slog.String("uri", uri),
					slog.String("error", e.Error()))

				continue
			}

			defer resp.Body.Close()
			_, _ = io.ReadAll(resp.Body)

			return

		case <-s.ctx.Done():
			return
		}
	}
}
