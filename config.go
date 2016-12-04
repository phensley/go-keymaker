package keymaker

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
)

// ClientConfig configures a drone client
type ClientConfig struct {

	// Dir is the directory for the config file, or cwd.
	Dir string

	// Addresses of drones in the cluster
	Addresses []string

	// BufferSize indicates number of keys to keep in the channel at all times
	BufferSize int `mapstructure:"buffer_size"`

	// Certificate file path containing the client's certificate in PEM
	Certificate string

	// PrivateKey file path containing the client's private key in PEM
	PrivateKey string `mapstructure:"private_key"`

	// CABundle file path containing the CA certificate bundle in PEM. Used to
	// authenticate drone certificates.
	CABundle string `mapstructure:"ca_bundle"`
}

// DroneConfig configures a Drone
type DroneConfig struct {

	// Dir is the directory for the config file, or cwd.
	Dir string

	// Address and port to listen on
	Address string

	// Concurrency level
	Concurrency int

	// Certificate file path containing the certificate in PEM
	Certificate string

	// PrivateKey file path containing the private key in PEM
	PrivateKey string `mapstructure:"private_key"`

	// CABundle file path containing the CA certificate bundle in PEM. Used to
	// authenticate client certificates.
	CABundle string `mapstructure:"ca_bundle"`

	// ClientAuth indicates strictness of client authentication
	ClientAuth string `mapstructure:"client_auth"`

	// ClientCN is a regular expression to match against the client
	// certificate subject's CommonName
	ClientCNRegexp string `mapstructure:"client_cn_regexp"`
}

// LoadConfigFile unmarshals the YAML contents of configPath
// and populates a config
func LoadConfigFile(cfg interface{}, configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return LoadConfig(cfg, data)
}

// LoadConfig unmarshals YAML and populates a config struct
func LoadConfig(cfg interface{}, config []byte) error {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		return err
	}
	return v.UnmarshalExact(cfg)
}

// LoadDroneTLSConfig ...
func LoadDroneTLSConfig(configDir string, cfg *DroneConfig) (*tls.Config, error) {
	certPEM, keyPEM, bundlePEM, err := loadCertificates(configDir, cfg.Certificate, cfg.PrivateKey, cfg.CABundle)
	if err != nil {
		return nil, err
	}
	return BuildDroneTLSConfig(certPEM, keyPEM, bundlePEM, cfg.ClientAuth)
}

// LoadClientTLSConfig loads in the x509 parts of the client's configuration
func LoadClientTLSConfig(configDir string, cfg *ClientConfig) (*tls.Config, error) {
	certPEM, keyPEM, bundlePEM, err := loadCertificates(configDir, cfg.Certificate, cfg.PrivateKey, cfg.CABundle)
	if err != nil {
		return nil, err
	}
	return BuildClientTLSConfig(certPEM, keyPEM, bundlePEM)
}

func loadCertificates(configDir string, certPath, keyPath, bundlePath string) ([]byte, []byte, []byte, error) {
	path := filepath.Join(configDir, certPath)
	certPEM, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error reading certificate from '%s': %s", path, err)
	}
	path = filepath.Join(configDir, keyPath)
	keyPEM, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error reading private key from '%s': %s", path, err)
	}
	path = filepath.Join(configDir, bundlePath)
	bundlePEM, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error reading ca bundle from '%s': %s", path, err)
	}
	return certPEM, keyPEM, bundlePEM, nil
}
