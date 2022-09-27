package tls

import (
	"crypto/tls"
	"fmt"
)

type (
	// IConfig is used to interface the TLS implemtation.
	IConfig interface {
		fmt.Stringer

		// GetCert return the full path to the server certificate file.
		GetCert() string

		// GetKey return the full path to the server key file.
		GetKey() string

		// GetCa return the full path to the server ca cert file.
		GetCa() string

		// GetInsecure return true if the TLS Certificate shouldn't be checked.
		GetInsecure() bool

		// GetLevel return
		GetLevel() tls.ClientAuthType

		// IsEmpty return true if the config is empty.
		Empty() bool
	}

	// Config contain the tls config passed by the config file.
	// It implement Config
	Config struct {
		Cert     string `json:"cert"     mapstructure:"cert"`
		Key      string `json:"key"      mapstructure:"key"`
		Ca       string `json:"ca"       mapstructure:"ca"`
		Insecure bool   `json:"insecure" mapstructure:"insecure"`
		Level    Level  `json:"level"    mapstructure:"level"`
	}
)

// GetCert implemte Config.
func (cfg Config) GetCert() string {
	return cfg.Cert
}

// GetKey implemte Config.
func (cfg Config) GetKey() string {
	return cfg.Key
}

// GetKey implemte Config.
func (cfg Config) GetCa() string {
	return cfg.Ca
}

// GetInsecure implemte Config.
func (cfg Config) GetInsecure() bool {
	return cfg.Insecure
}

func (cfg Config) GetLevel() tls.ClientAuthType {
	return cfg.Level.STD()
}

// Empty implemte Config.
func (cfg Config) Empty() bool {
	return cfg.Cert == "" && cfg.Key == ""
}

// String implement Stringer interface.
func (cfg Config) String() string {
	if cfg.Empty() {
		return ""
	}

	return fmt.Sprintf("\t ~!~ cert:\t%q\n\t ~!~ key:\t%q\n\t ~!~ ca:\t%q,\n\t ~!~ insecure:\t%t\n\t ~!~ level:\t%s\n",
		cfg.Cert, cfg.Key, cfg.Ca, cfg.Insecure, cfg.Level)
}
