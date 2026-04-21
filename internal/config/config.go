package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	AppPort     string `mapstructure:"APP_PORT"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	RedisHost   string `mapstructure:"REDIS_HOST"`
	RedisPort   string `mapstructure:"REDIS_PORT"`
	OTELCollector string `mapstructure:"OTEL_COLLECTOR"`
}

func LoadConfig() *Config {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}

	return &config
}
