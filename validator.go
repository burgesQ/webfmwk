package webfmwk

import (
	en_translator "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type (
	// ValidationError is returned in case of form / query validation error
	ValidationError struct {
		Status int                                    `json:"status"`
		Error  validator.ValidationErrorsTranslations `json:"message"`
	}
)

var (
	// validate annotation : `validate` : go-playground
	validate *validator.Validate
	// universal translator
	uni *ut.UniversalTranslator

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans ut.Translator
)

func initValidator() {
	var en = en_translator.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en") // we know that en exist

	validate = validator.New()
	if e := en_translations.RegisterDefaultTranslations(validate, trans); e != nil {
		logger.Fatalf("cannot init translations : %s", e.Error())
	}
}

func GetValidator() *validator.Validate {
	// from init server - if validator is called before
	// the server init (which happend pretty often)
	once.Do(initOnce)
	return validate
}

func RegisterValidatorAlias(name, what string) {
	// from init server - if validator is called before
	// the server init (which happend pretty often)
	once.Do(initOnce)
	validate.RegisterAlias(name, what)
}
