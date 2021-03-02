package security

import "github.com/burgesQ/webfmwk/v5"

const (
	HeaderProtection = "X-Xss-Protection"
	HeaderOption     = "X-Content-Type-Options"
	HeaderSecu       = "Strict-Transport-Security"

	HeaderProtectionV = "1; mode=block"
	HeaderOptionV     = "nosniff"
	HeaderSecuV       = "max-age=3600; includesubDomains"
)

// Handler append few security headers
func Handler(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		c.SetHeaders(webfmwk.Header{HeaderProtection, HeaderProtectionV},
			webfmwk.Header{HeaderOption, HeaderOptionV},
			webfmwk.Header{HeaderSecu, HeaderSecuV})

		return next(c)
	})
}
