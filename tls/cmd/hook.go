package cmd

import (
	"reflect"

	"github.com/burgesQ/webfmwk/v5/tls"
	"github.com/mitchellh/mapstructure"
)

// StringToLevelHookFunc allow tls.Level param to be processed
// by the spf13/cobra and spf13/viper utility.
// Use as following with frafos:
//
//	func fetchCfg() (cfg api.ServerCfg) {
//		cmd.ReadCfg(&cfg,
//			cmd_log.StringToSyslogHookFunc(),
//			tls.StringToLevelHookFunc())
func StringToLevelHookFunc() mapstructure.DecodeHookFunc {
	return stringToLevelHookFunc
}

func stringToLevelHookFunc(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Name() != "string" || t.Name() != "Level" {
		return data, nil
	}

	raw, ok := data.(string)
	if raw == "" || !ok {
		return data, nil
	} else if t != reflect.TypeOf(tls.NoClientCert) {
		return data, nil
	}

	var lvl tls.Level
	e := lvl.Set(raw)

	return lvl, e
}
