package tls

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// StringToSyslogHookFunc allow tls.Level param to be processed
// by the spf13/cobra and spf13/viper utility.
// Use as following with frafos:
//
//  func fetchCfg() (cfg api.ServerCfg) {
//          cmd.ReadCfg(&cfg,
//                  cmd_log.StringToSyslogHookFunc(),
//                  tls.StringToSyslogHookFunc())
func StringToSyslogHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (
		interface{}, error) {
		if f.Name() != "string" || t.Name() != "Level" {
			return data, nil
		}

		raw, ok := data.(string)
		if raw == "" || !ok {
			return data, nil
		} else if t != reflect.TypeOf(NoClientCert) {
			return data, nil
		}

		lvl := NoClientCert
		e := lvl.Set(raw)

		return lvl, e
	}
}
