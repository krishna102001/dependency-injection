package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// create a central config where all other different config will combine
type Config struct {
	Server *ServerConfig `yaml:"server"`
	Mongo  *MongoConfig  `yaml:"mongo"`
}

// server config
type ServerConfig struct {
	Port      string `yaml:"port"`
	JwtSecret string `yaml:"jwtSecret"`
}

// database config
type MongoConfig struct {
	URI        string `yaml:"uri"`
	DBName     string `yaml:"dbName"`
	Collection string `yaml:"collection"`
}

// config going to store the since its small because i don't want to expose it to other packages
var config *Config

// it will going to check config is weather loaded or not
var isConfigLoaded bool

// get config will be called when we need to access the config data
func GetConfig() (*Config, error) {
	if !isConfigLoaded {
		return nil, fmt.Errorf("config has not been loaded in main.go")
	}
	return config, nil
}

// it will load the config from config.yaml file
func LoadConfigLocal() (*Config, error) {
	// os.Getwd() will give you current directory for me it will give dependecy-injection
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// we will join cwd with config folder
	configDir := filepath.Join(cwd, "config")

	// pass to the centralConfig because we can have one config for Testing also
	// thats why we have separated load of config
	cfg, err := LoadCentralConfig(configDir)
	if err != nil {
		isConfigLoaded = false
		return nil, err
	}
	isConfigLoaded = true
	config = cfg
	return config, nil
}

// this will call generic function to load the config file
func LoadCentralConfig(configDir string) (*Config, error) {
	readConfig, err := loadYaml[Config](filepath.Join(configDir, "config.yaml"))
	if err != nil {
		return nil, err
	}
	return readConfig, nil
}

// it read the file and unmarshal the yaml file and assign to config
func loadYaml[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file(%s): %v", filename, err)
	}
	// it will create a new T type with pointer(var config *T)
	config := new(T)
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the yaml file %v", err)
	}
	return config, nil
}
