package tls

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
)

type Level tls.ClientAuthType

const (
	// copied from https://pkg.go.dev/crypto/tls#ClientAuthType

	// NoClientCert indicates that no client certificate should be requested
	// during the handshake, and if any certificates are sent they will not
	// be verified.
	NoClientCert Level = iota
	// RequestClientCert indicates that a client certificate should be requested
	// during the handshake, but does not require that the client send any
	// certificates.
	RequestClientCert
	// RequireAnyClientCert indicates that a client certificate should be requested
	// during the handshake, and that at least one certificate is required to be
	// sent by the client, but that certificate is not required to be valid.
	RequireAnyClientCert
	// VerifyClientCertIfGiven indicates that a client certificate should be requested
	// during the handshake, but does not require that the client sends a
	// certificate. If the client does send a certificate it is required to be
	// valid.
	VerifyClientCertIfGiven
	// RequireAndVerifyClientCert indicates that a client certificate should be requested
	// during the handshake, and that at least one valid certificate is required
	// to be sent by the client.
	RequireAndVerifyClientCert
	// RequireAndVerifyClientCertAndSAN is the same as RequireAndVerifyClientCert
	// with an extra check to the certificate SAN.
	RequireAndVerifyClientCertAndSAN
)

var (
	_lvl2str = map[Level]string{
		NoClientCert:                     "never",
		RequestClientCert:                "demande",
		RequireAnyClientCert:             "allow",
		VerifyClientCertIfGiven:          "try",
		RequireAndVerifyClientCert:       "hard",
		RequireAndVerifyClientCertAndSAN: "hardAndSAN",
	}

	_str2lvl = map[string]Level{
		"never":      NoClientCert,
		"demande":    RequestClientCert,
		"allow":      RequireAnyClientCert,
		"try":        VerifyClientCertIfGiven,
		"hard":       RequireAndVerifyClientCert,
		"hardAndSAN": RequireAndVerifyClientCertAndSAN,
	}

	_toNatif = map[Level]tls.ClientAuthType{
		NoClientCert:                     tls.NoClientCert,
		RequestClientCert:                tls.RequestClientCert,
		RequireAnyClientCert:             tls.RequireAnyClientCert,
		VerifyClientCertIfGiven:          tls.VerifyClientCertIfGiven,
		RequireAndVerifyClientCert:       tls.RequireAndVerifyClientCert,
		RequireAndVerifyClientCertAndSAN: tls.RequireAndVerifyClientCert,
	}
)

type LevelError struct{ what string }

func (e LevelError) Error() string {
	return fmt.Sprintf("non existing TLS level: %s", e.what)
}

func (lv Level) String() string {
	if v, ok := _lvl2str[lv]; ok {
		return v
	}

	return `{}`
}

func (lv Level) STD() tls.ClientAuthType {
	return _toNatif[lv]
}

func (lv *Level) Set(val string) error {
	if v, ok := _str2lvl[val]; ok {
		*lv = v

		return nil
	}

	return LevelError{val}
}

func (lv *Level) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	return lv.Set(j)
}

func (lv Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(lv.String())
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}
