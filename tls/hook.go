package tls

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

const _Str = "string"

func StringToSyslogHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		fmt.Printf("handling -> %q - %q\n", f.Name(), t.Name())

		if f.Name() != "string" || t.Name() != "Level" {
			return data, nil
		}

		raw, ok := data.(string)
		if raw == "" || !ok {
			return data, nil
		}

		fmt.Printf("raw -> %q\n", raw)

		switch t {
		case reflect.TypeOf(NoClientCert):
			lvl := NoClientCert
			e := lvl.Set(raw)

			return lvl, e
		}

		return data, nil
	}
}
