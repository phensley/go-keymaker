package keymaker

import (
	"bytes"
	"io/ioutil"

	"github.com/spf13/viper"
)

// LoadConfigFile unmarshals the YAML defaults and contents of configPath
// and populates the given config struct
func LoadConfigFile(cfg interface{}, configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	return LoadConfig(cfg, data)
}

// LoadConfig unmarshals the YAML defaults and config and populates
// the given config struct
func LoadConfig(cfg interface{}, config []byte) error {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		return err
	}
	return v.UnmarshalExact(cfg)
}
