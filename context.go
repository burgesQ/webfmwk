package webfmwk

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/valyala/fasthttp"
)

const _prettyTag = "pretty"

type (
	// InputHandling interface introduce I/O actions.
	InputHandling interface {
		// FetchContent extract the json content from the body into the content interface.
		FetchContent(content interface{}) ErrorHandled

		// Validate is used to validate a content of the content params.
		// See https://go-playground/validator.v10 for more.
		Validate(content interface{}) ErrorHandled

		// FetchAndValidateContent fetch the content then validate it.
		FetchAndValidateContent(content interface{}) ErrorHandled

		// DecodeQP load the query param in the content object.
		// Seee https://github.com/gorilla/query for more.
		DecodeQP(content interface{}) ErrorHandled

		// DecodeAndValidateQP load the query param in the content object and then validate it.
		DecodeAndValidateQP(content interface{}) ErrorHandled
	}

	// ContextLogger interface implement the context Logger needs.
	ContextLogger interface {
		// SetStructuredLogger set context' structured logger.
		SetStructuredLogger(logger *slog.Logger) Context

		// GetStructuredLogger return context' structured logger.
		GetStructuredLogger() *slog.Logger
	}

	// Context interface implement the context used in this project.
	Context interface {
		SendResponse
		InputHandling
		ContextLogger

		// GetFastContext return a pointer to the internal fasthttp.RequestCtx.
		GetFastContext() *fasthttp.RequestCtx

		// GetContext return the request context.Context.
		GetContext() context.Context

		// GetVar return the url var parameters. An empty string for missing case.
		GetVar(key string) (val string)

		// GetQueries return the queries into a fasthttp.Args object.
		GetQuery() *fasthttp.Args

		// GetQuery fetch the query object key
		// GetQuery(key string) (val string, ok bool)
	}

	// icontext implement the Context interface
	// It hold the data used by the request
	icontext struct {
		*fasthttp.RequestCtx
		slog *slog.Logger
		ctx  context.Context //nolint:containedctx
	}
)

var (
	// decoder annotation : `schema` : gorilla
	decoder                 = schema.NewDecoder()
	errUnprocessablePayload = NewUnprocessable(NewError("Unprocessable payload"))
)

// GetVar implement Context
func (c *icontext) GetVar(key string) string {
	v, ok := c.UserValue(key).(string)
	if !ok {
		return ""
	}

	return v
}

// GetQuery implement Context.
func (c *icontext) GetQuery() *fasthttp.Args {
	return c.QueryArgs()
}

// GetFastContext implement Context.
func (c *icontext) GetFastContext() *fasthttp.RequestCtx {
	return c.RequestCtx
}

// SetStructuredLogger implement Context.
func (c *icontext) SetStructuredLogger(logger *slog.Logger) Context {
	c.slog = logger

	return c
}

// GetLogger implement Context
func (c *icontext) GetStructuredLogger() *slog.Logger {
	return c.slog
}

// GetContext implement Context
func (c *icontext) GetContext() context.Context {
	return c.ctx
}

// FetchContent implement Context.
// It load payload in the dest interface{} using the system json library.
func (c *icontext) FetchContent(dest interface{}) ErrorHandled {
	b := c.PostBody()

	if e := json.Unmarshal(b, &dest); e != nil {
		c.slog.Error("fetching payload", slog.Any("error", e))

		return errUnprocessablePayload
	}

	return nil
}

// Validate implement Context
// this implemtation use validator to anotate & check struct
func (c *icontext) Validate(dest interface{}) ErrorHandled {
	if e := validate.Struct(dest); e != nil {
		c.slog.Error("validating form or query param", slog.Any("error", e))

		var ev validator.ValidationErrors

		errors.As(e, &ev)

		return NewUnprocessable(ValidationError{
			Status: http.StatusUnprocessableEntity,
			Error:  TranslateAndUseFieldName(ev),
		})
	}

	return nil
}

// FetchAndValidateContent implemt Context.
// It sucesively call FetchContent then Validate on the dest param
func (c *icontext) FetchAndValidateContent(dest interface{}) ErrorHandled {
	if e := c.FetchContent(&dest); e != nil {
		return e
	}

	return c.Validate(dest)
}

// DecodeQP implement Context
func (c *icontext) DecodeQP(dest interface{}) ErrorHandled {
	m := map[string][]string{}

	c.QueryArgs().VisitAll(func(k, v []byte) {
		key := string(k)
		m[key] = make([]string, 1)
		m[key][0] = string(v)
	})

	if e := decoder.Decode(dest, m); e != nil {
		c.slog.Error("validating query params", slog.Any("error", e))

		return NewUnprocessable(NewErrorFromError(e))
	}

	return nil
}

// DecodeAndValidateQP implement Context
func (c *icontext) DecodeAndValidateQP(qp interface{}) ErrorHandled {
	if e := c.DecodeQP(qp); e != nil {
		return e
	}

	return c.Validate(qp)
}
