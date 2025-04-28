package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *Server
		Db       *Db
		Session  *Session
		Csrf     *Csrf
		RabbitMq *RabbitMq `mapstructure:"rabbitmq"`
	}

	Server struct {
		Port int
		Name string
		Url  string
		Os   string
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}

	Session struct {
		Name       string `mapstructure:"name"`
		Secret     string `mapstructure:"secret"`
		Expiration int    `mapstructure:"expiration"`
	}

	Csrf struct {
		Secret     string `mapstructure:"secret"`
		Name       string `mapstructure:"name"`
		Expiration int    `mapstructure:"expiration"`
	}

	RabbitMq struct {
		Host  string `mapstructure:"host"`
		Queue string `mapstructure:"queue"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
	})

	return configInstance
}
