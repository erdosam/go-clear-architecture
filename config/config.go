package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		JwtAuth `yaml:"jwtauth"`
		Casbin  `yaml:"casbin"`
		//RMQ  `yaml:"rabbitmq"`
		Google `yaml:"google"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port    string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		AuthKey string `env-required:"true" yaml:"auth_key" env:"HTTP_AUTH_KEY"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// JwtAuth -.
	JwtAuth struct {
		BaseUrl       string `env-required:"true" yaml:"base_url" env:"JWTAUTH_BASE_URL"`
		ClientKeyFile string `env-required:"true" yaml:"client_key_file" env:"JWTAUTH_CLIENT_KEY_FILE"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true"                            env:"RMQ_URL"`
	}

	Casbin struct {
		ModelFile          string `env-required:"true" yaml:"model_file" env:"CASBIN_MODEL_FILE"`
		LoadPolicyInterval int    `env-required:"true" yaml:"load_policy_interval" env:"LOAD_POLICY_INTERVAL"`
	}

	Google struct {
		ProjectId  string `env-required:"true" yaml:"project_id" env:"GOOGLE_PROJECT_ID"`
		BucketName string `env-required:"true" yaml:"bucket_name" env:"GOOGLE_BUCKET_NAME"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
