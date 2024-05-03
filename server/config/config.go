package config

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"server/utils"
)

//go:embed config.yaml
var DEFAULT_CONFIG_FILE string

var config *Config

type Config struct {
	API struct {
		Address      string   `yaml:"address"`
		AllowOrigins []string `yaml:"allow_origins"`
		MemberFile   string   `yaml:"member_file"`
		JwtExpireDay int      `yaml:"jwt_expire_day"`
	} `yaml:"api"`
	Cache struct {
		Root     string `yaml:"root"`
		UUIDFile string `yaml:"uuid_file"`
	} `yaml:"cache"`
	Log struct {
		Root      string `yaml:"root"`
		WithColor bool   `yaml:"with_color"`
	} `yaml:"log"`
}

func LoadConfig(path string, config *Config) (err error) {
	if !utils.HasFile(path) {
		if err = utils.AutoWriteFile(path, []byte(DEFAULT_CONFIG_FILE), os.ModePerm); err != nil {
			return fmt.Errorf("failed to write default config file: %w", err)
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return yaml.Unmarshal(content, config)
}

func (config Config) GetCacheUUIDPath() string {
	return path.Join(config.Cache.Root, config.Cache.UUIDFile)
}

func SetupGlobalConfig(path string) error {
	if config == nil {
		config = &Config{}
	}
	if err := LoadConfig(path, config); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

func Get() *Config {
	return config
}
