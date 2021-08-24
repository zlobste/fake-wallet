package config

import (
	"fmt"
	_ "github.com/lib/pq" // here
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config interface {
	Listener() string
	JWTSecret() string

	Logger
	Databaser
}

type config struct {
	Addr         string `yaml:"addr"`
	Log          string `yaml:"log"`
	DatabaseUrl  string `yaml:"db_url"`
	JWTSecretKey string `yaml:"jwt_secret_key"`

	Logger
	Databaser
}

func New(path string) Config {
	cfg := config{}

	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to read config: %s", path)))
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to unmarshal config: %s", path)))
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.Addr = ":" + port
	}
	cfg.Logger = NewLogger(cfg.Log)
	cfg.Databaser = NewDatabaser(cfg.DatabaseUrl, cfg.Logger.Logging())

	return &cfg
}

func (c *config) Listener() string {
	return c.Addr
}

func (c *config) JWTSecret() string {
	return c.JWTSecretKey
}
