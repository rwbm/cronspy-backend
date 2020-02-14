package config

import (
	"fmt"
	"io/ioutil"

	"github.com/astropay/go-tools/files"
	yaml "gopkg.in/yaml.v2"
)

// Configuration is the structure used to hold configuration from config.yml
type Configuration struct {
	Server struct {
		Port            string `yaml:"port"`
		ReadTimeout     int    `yaml:"read_timeout"`
		WriteTimeout    int    `yaml:"write_timeout"`
		Debug           bool   `yaml:"debug"`
		TokenExpiration int    `yaml:"token_expiration"`
	} `yaml:"server"`
	Database struct {
		Driver             string `yaml:"driver"`
		Address            string `yaml:"address"`
		DefaultDB          string `yaml:"default_db"`
		Username           string `yaml:"username"`
		Password           string `yaml:"password"`
		MaxOpenConnections int    `yaml:"max_open_connections"`
		MaxIdleConnections int    `yaml:"max_idle_connections"`
		MaxLifeTime        int    `yaml:"max_lifetime"`
	} `yaml:"database"`
}

// Load reads application settings in the indicated file
func Load(path string) (cfg *Configuration, err error) {
	if files.Exists(path) {
		cfg, err = loadSettingsFromFile(path)
	} else {
		return nil, fmt.Errorf("file '%s' not found", path)
	}

	return
}

// load settings from yaml file
func loadSettingsFromFile(path string) (cfg *Configuration, err error) {
	cfg = new(Configuration)

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %s", err)
	}

	if err := yaml.Unmarshal(fileContent, cfg); err != nil {
		return nil, fmt.Errorf("error parsing configuration file: %s", err)
	}

	return
}
