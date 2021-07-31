package webfmwk

import (
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
		Status int                                    `json:"status"`
		Error  validator.ValidationErrorsTranslations `json:"message"`
	}
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

func initValidator() {
	var (
		en = en_translator.New()
		ok bool
	)

	uni = ut.New(en, en)
	if trans, ok = uni.GetTranslator("en"); !ok {
		logger.Fatalf("cannot get en translator")
	}

	validate = validator.New()
	if e := en_translations.RegisterDefaultTranslations(validate, trans); e != nil {
		logger.Fatalf("cannot init translations : %s", e.Error())
	}

	useJSONFieldName()
}

// Use the struct json field name for validation errors
// src: https://github.com/go-playground/validator/blob/9a5bce32538f319bf69aebb3aca90d394bc6d0cb/_examples/struct-level/main.go#L37
func useJSONFieldName() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}

// RegisterValidatorRule register the  validation rule param.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorRule(name string, fn func(fl validator.FieldLevel) bool) error {
	once.Do(initOnce)
	return validate.RegisterValidation(name, fn)
}

// RegisterValidatorAlias register some validation alias.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorAlias(name, what string) {
	// from init server - if validator is called before
	// the server init (which may happen pretty often)
	once.Do(initOnce)
	validate.RegisterAlias(name, what)
}

// RegisterValidatorTrans register some validation alias.
// See https://go-playground/validator.v10 for more.
func RegisterValidatorTrans(name, what string) error {
	return validate.RegisterTranslation(name, trans,
		func(ut ut.Translator) error {
			return ut.Add(name, what, true) // see universal-translator for details
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(name, fe.Field())
			return t
		})
}
