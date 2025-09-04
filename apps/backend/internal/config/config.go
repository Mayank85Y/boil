package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct{
	Primary			Primary					`koanf:"primary" validate:"required"`
	Server			ServerConfig			`koanf:"server" validate:"required"`
	Database		DatabaseConfig			`koanf:"database" validate:"required"`
	Auth			AuthConfig				`koanf:"auth" validate:"required"`
	Redis			RedisConfig				`koanf:"redis" validate:"required"`
	Integration 	IntegrationConfig		`koanf:"integration" validate:"required"`
	Observability	*ObservabilityConfig 	`koanf:"observability"`
}

type Primary struct{
	Env string `koanf:"env" validate:"required"` //`` is go struct tags heelp in reflection help some metadaata
}

type ServerConfig struct {
	Port				string	 `koanf:"port" validate:"required"`	
	ReadTimeout			int		 `koanf:"read_timeout" validate:"required"`
	WriteTimeout		int		 `koanf:"write_timeout" validate:"required"`
	IdleTimeout			int		 `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins	[]string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host			string	`koanf:"host" validated:"requred"`
	Port			int		`koanf:"port" validated:"requred"`
	User			string	`koanf:"user" validated:"requred"`
	Password		string	`koanf:"password"`
	Name			string	`koanf:"name" validated:"requred"`
	SSLMode			string	`koanf:"ssl_mode" validated:"requred"`
	MaxOpenConns	int		`koanf:"max_open_conns" validated:"requred"`
	MaxIdleConns   	int		`koanf:"max_idle_conns" validated:"requred"`
	ConnMaxLifetime	int		`koanf:"conn_max_life_time" validated:"requred"`
	ConnMaxIdleTime	int		`koanf:"conn_max_idle_time" validated:"requred"`
}

type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validated:"required"`
}

type RedisConfig struct {
	Address	string	`koanf:"address" validate:"required"`
}

type IntegrationConfig struct{
	ResendAPIKey string `koanf:"resend_api_key" validate:"required"`
}

func LoadConfig() (*Config, error){
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")
	err := k.Load(env.Provider("BOIL_", ".", func(s string) string{
		return strings.ToLower(strings.TrimPrefix(s, "BOIL_")) 
	}), nil)
	if err != nil{
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}
	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil{
		logger.Fatal().Err(err).Msg("Could not unmarshal main config")
	}

	validate := validator.New()
	err = validate.Struct(mainConfig)
	if err != nil{
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	mainConfig.Observability.ServiceName = "boil"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("Invalid observability config")
	}
	return mainConfig, nil
}