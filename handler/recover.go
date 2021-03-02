package handler

import . "github.com/burgesQ/webfmwk/v4"

// Recover launch a panic catcher - if the catched panic hold an
// webfmwk.ErrorHandled then a API error response is generated from it.
func Recover(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case ErrorHandled:
					_ = c.JSON(e.GetOPCode(), e.GetContent())
				default:
					c.GetLogger().Errorf("catched %T %#v", e, e)
					panic(e)
				}
			}
		}()

		return next(c)
	})
}
