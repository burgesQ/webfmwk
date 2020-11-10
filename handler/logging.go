package handler

import "github.com/burgesQ/webfmwk/v4"

// Logging log information about the newly receive request
func Logging(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		r := c.GetRequest()

		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		c.GetLogger().Infof("[+] (%s) %s : [%s]%s", c.GetRequestID(), IPAddress, r.Method, r.RequestURI)

		return next(c)
	})
}
