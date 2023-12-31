package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Database   database   `yaml:"database"`
		FusionAuth fusionAuth `yaml:"fusion_auth"`
		Server     server     `yaml:"server"`
		LogLevel   string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"dev"`
		Flags      flags      `yaml:"flags"`
	}

	server struct {
		Port         string        `yaml:"port" env:"PORT" env-default:":8080"`
		TimeoutRead  time.Duration `yaml:"timeout_read" env:"TIMEOUT_READ" env-default:"5s"`
		TimeoutWrite time.Duration `yaml:"timeout_write" env:"TIMEOUT_WRITE" env-default:"5s"`
	}

	fusionAuth struct {
		Host   string `yaml:"host" env:"FUSION_AUTH_HOST" env-default:"http://localhost:9011"`
		AppId  string `yaml:"app_id" env:"FUSION_AUTH_APP_ID" env-default:"poc-auth"`
		ApiKey string `yaml:"api_key" env:"FUSION_AUTH_API_KEY" env-required:"true"`
	}

	database struct {
		MongoDB mongoConfig `yaml:"mongodb"`
	}

	mongoConfig struct {
		Host     string `yaml:"host" env:"MONGO_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"MONGO_PORT" env-default:"27017"`
		Database string `yaml:"database" env:"MONGO_DATABASE" env-default:"auth"`
		Username string `yaml:"username" env:"MONGO_USERNAME" env-default:"root"`
		Password string `yaml:"password" env:"MONGO_PASSWORD" env-default:"root"`
		URI      string `yaml:"uri" env:"MONGO_URI"`
	}

	flags struct {
		ConfigPath string
		Migrations bool
	}
)

func (mondodb mongoConfig) GetURI() string {
	if mondodb.URI != "" {
		return mondodb.URI
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		mondodb.Username,
		mondodb.Password,
		mondodb.Host,
		mondodb.Port,
		mondodb.Database,
	)
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Flags: loadFlags(),
	}

	if cfg.Flags.ConfigPath != "" {
		err := cleanenv.ReadConfig(cfg.Flags.ConfigPath, &cfg)
		if err != nil {
			return cfg, fmt.Errorf("failed to read config file: %w", err)
		}

		return cfg, nil
	} else {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return cfg, fmt.Errorf("failed to read env: %w", err)
		}
	}

	return cfg, nil
}

func loadFlags() flags {
	f := flags{}
	flag.StringVar(&f.ConfigPath, "config", "", "path to config file")
	flag.BoolVar(&f.Migrations, "migrations", false, "run migrations")
	flag.Parse()
	return f
}
