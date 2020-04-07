package handler

import "github.com/burgesQ/webfmwk/v4"

// Logging log information about the newly receive request
func Logging(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.IContext) {
		r := c.GetRequest()
		c.GetLogger().Infof("[+] (%s) : [%s]%s", c.GetRequestID(), r.Method, r.RequestURI)
		next(c)
	})
}
