package tls

import "fmt"

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

		// IsEmpty return true if the config is empty.
		Empty() bool
	}

	// Config contain the tls config passed by the config file.
	// It implement Config
	Config struct {
		Cert     string `json:"cert" mapstructur:"cert"`
		Key      string `json:"key" mapstructur:"key"`
		Ca       string `json:"ca" mapstructur:"ca"`
		Insecure bool   `json:"insecure" mapstructur:"insecure"`
	}
)

// GetCert implemte Config.
func (config Config) GetCert() string {
	return config.Cert
}

// GetKey implemte Config.
func (config Config) GetKey() string {
	return config.Key
}

// GetKey implemte Config.
func (config Config) GetCa() string {
	return config.Ca
}

// GetInsecure implemte Config.
func (config Config) GetInsecure() bool {
	return config.Insecure
}

// Empty implemte Config.
func (config Config) Empty() bool {
	return config.Cert == "" && config.Key == ""
}

// String implement Stringer interface.
func (config Config) String() string {
	if config.Empty() {
		return ""
	}

	return fmt.Sprintf("\t ~!~ cert:\t%q\n\t ~!~ key:\t%q\n\t ~!~ ca:\t%q,\n\t ~!~ insecure:\t%t\n",
		config.Cert, config.Key, config.Ca, config.Insecure)
}
