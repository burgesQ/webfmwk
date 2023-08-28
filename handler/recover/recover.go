//nolint:predeclared
package recover

import (
	"fmt"

	"github.com/burgesQ/webfmwk/v6"
)

// Handler launch a panic catcher - if the catched panic hold an
// webfmwk.ErrorHandled then a API error response is generated from it.
func Handler(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case error:
					webfmwk.HandleError(c, e)

				default:
					c.GetStructuredLogger().Error("catched exit", "error", e)
					_ = c.JSONInternalError(webfmwk.NewError(
						fmt.Sprintf("internal error: %T %v", e, e)))
				}
			}
		}()

		return next(c)
	})
}
