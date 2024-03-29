package webfmwk

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	en_translator "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type (
	// ValidationError is returned in case of form / query validation error
	// see gtihub.com/go-playground/validator.v10
	ValidationError struct {
		Error  ErrorValidation `json:"message"`
		Status int             `json:"status"`
	}

	// ErrorsValidation is a map of translated errors
	ErrorValidation map[string]string
)

var (
	//  validate annotation : `validate` : go-playground
	validate *validator.Validate
	// universal translator
	uni *ut.UniversalTranslator

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans ut.Translator
)

var ErrGetTranslator = errors.New("fetching the 'en' translator")

func initValidator() error {
	var (
		en = en_translator.New()
		ok bool
	)

	uni = ut.New(en, en)
	if trans, ok = uni.GetTranslator("en"); !ok {
		return ErrGetTranslator
	}

	validate = validator.New()
	if e := en_translations.RegisterDefaultTranslations(validate, trans); e != nil {
		return fmt.Errorf("initializing the 'en' translator : %w", e)
	}

	useJSONFieldName()

	return nil
}

// Trnaslate the errs array of validation error and use the actual filed name
// instad of the full struct namepsace one.
// src: https://blog.depa.do/post/gin-validation-errors-handling#toc_8
func TranslateAndUseFieldName(errs validator.ValidationErrors) ErrorValidation {
	es := ErrorValidation{}

	for _, f := range errs {
		es[f.Field()] = f.Translate(trans)
	}

	return es
}

// Use the struct json field name for validation errors
// src: shorturl.at/evDJ0
func useJSONFieldName() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		//nolint: gomnd
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}

// RegisterValidatorRule register the  validation rule param.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorRule(name string, fn func(fl validator.FieldLevel) bool) (e error) {
	once.Do(func() { e = initOnce() })

	if e != nil {
		return fmt.Errorf("registering validator rule: %w", e)
	}

	return validate.RegisterValidation(name, fn)
}

// RegisterValidatorAlias register some validation alias.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorAlias(name, what string) (e error) {
	// from init server - if validator is called before
	// the server init (which may happen pretty often)
	once.Do(func() { e = initOnce() })

	if e != nil {
		return fmt.Errorf("registering validator alias: %w", e)
	}

	validate.RegisterAlias(name, what)

	return nil
}

// RegisterValidatorTrans register some validation alias.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorTrans(name, what string) (e error) {
	// from init server - if validator is called before
	// the server init (which may happen pretty often)
	once.Do(func() { e = initOnce() })

	if e != nil {
		return fmt.Errorf("init validator: %w", e)
	}

	return validate.RegisterTranslation(name, trans,
		func(ut ut.Translator) error {
			return ut.Add(name, what, true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(name, fe.Field())

			return t
		})
}
