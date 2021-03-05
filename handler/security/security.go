package security

import "github.com/burgesQ/webfmwk/v5"

const (
	headerProtection = "X-Xss-Protection"
	headerOption     = "X-Content-Type-Options"
	headerSecu       = "Strict-Transport-Security"

	headerProtectionV = "1; mode=block"
	headerOptionV     = "nosniff"
	headerSecuV       = "max-age=3600; includesubDomains"
)

// Handler append few security headers
func Handler(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		c.SetHeaders(webfmwk.Header{headerProtection, headerProtectionV},
			webfmwk.Header{headerOption, headerOptionV},
			webfmwk.Header{headerSecu, headerSecuV})

		return next(c)
	})
}
