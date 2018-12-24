package web

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/exception"
)

// TLSConfig is a config for app tls settings.
type TLSConfig struct {
	Cert     []byte `json:"cert,omitempty" yaml:"cert,omitempty" env:"TLS_CERT"`
	CertPath string `json:"certPath,omitempty" yaml:"certPath,omitempty" env:"TLS_CERT_PATH"`
	Key      []byte `json:"key,omitempty" yaml:"key,omitempty" env:"TLS_KEY"`
	KeyPath  string `json:"keyPath,omitempty" yaml:"keyPath,omitempty" env:"TLS_KEY_PATH"`

	CAPaths []string `json:"caPaths,omitempty" yaml:"caPaths,omitempty" env:"TLS_CA_PATHS,csv"`
}

// GetCert returns a tls cert.
func (tc TLSConfig) GetCert(defaults ...[]byte) []byte {
	return configutil.CoalesceBytes(tc.Cert, nil, defaults...)
}

// GetCertPath returns a tls cert path.
func (tc TLSConfig) GetCertPath(defaults ...string) string {
	return configutil.CoalesceString(tc.CertPath, "", defaults...)
}

// GetKey returns a tls key.
func (tc TLSConfig) GetKey(defaults ...[]byte) []byte {
	return configutil.CoalesceBytes(tc.Key, nil, defaults...)
}

// GetKeyPath returns a tls key path.
func (tc TLSConfig) GetKeyPath(defaults ...string) string {
	return configutil.CoalesceString(tc.KeyPath, "", defaults...)
}

// GetCAPaths returns a list of ca paths to add.
func (tc TLSConfig) GetCAPaths(defaults ...[]string) []string {
	return configutil.CoalesceStrings(tc.CAPaths, nil, defaults...)
}

// GetConfig returns a stdlib tls config for the config.
func (tc TLSConfig) GetConfig() (*tls.Config, error) {
	if !tc.HasKeyPair() {
		return nil, nil
	}

	var cert tls.Certificate
	var err error

	if len(tc.GetCertPath()) > 0 {
		cert, err = tls.LoadX509KeyPair(
			tc.GetCertPath(),
			tc.GetKeyPath(),
		)
	} else {
		cert, err = tls.X509KeyPair(tc.GetCert(), tc.GetKey())
	}

	if err != nil {
		return nil, exception.New(err)
	}

	if len(tc.GetCAPaths()) == 0 {
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
		}, nil
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, exception.New(err)
	}
	for _, caPath := range tc.GetCAPaths() {
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, exception.New(err)
		}
		certPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS11,
	}, nil
}

// HasKeyPair returns if the config names a keypair.
func (tc TLSConfig) HasKeyPair() bool {
	if len(tc.GetCert()) > 0 && len(tc.GetKey()) > 0 {
		return true
	}

	if len(tc.GetCertPath()) > 0 && len(tc.GetKeyPath()) > 0 {
		return true
	}

	return false
}
