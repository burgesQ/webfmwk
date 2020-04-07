// Package handler implement some extra handler to the webfmwk.
// Handler provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which may perform various
// operations and should finish by calling the next HTTP handler.
// Middleware function signature is `func(next webfmwk.HandlerFunc) webfmwk.HandlerFunc` an
// webfmwk.HandlerFunc is `func(c webfmwk.IContext)`
package handler
