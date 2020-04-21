package handler

import "github.com/burgesQ/webfmwk/v4"

// Recover launch a panic catcher. If the catched panic hold an
// webfmwk.ErrorHandled then a response is generated from it
func Recover(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case webfmwk.ErrorHandled:
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
