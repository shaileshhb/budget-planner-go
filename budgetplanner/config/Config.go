package config

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/spf13/viper"
)

// Config Contain Viper
type Config struct {
	viper *viper.Viper
}

// ConfReader defines all methods to be present in Config.
type ConfReader interface {
	GetString(key EnvKey) string
	IsSet(key EnvKey) bool
	GetInt64(key EnvKey) int64
}

// NewConfig Read envfile and Return Config
func NewConfig(isProduction bool) ConfReader {

	vp := viper.New()
	if isProduction {
		vp.SetConfigName("config")
	} else {
		vp.SetConfigName("dev-config")
	}
	vp.SetConfigType("env")
	vp.AddConfigPath(".")
	vp.AutomaticEnv()

	config := Config{
		viper: vp,
	}

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.GetLogger().Warn("file Not Found")
		} else {
			log.GetLogger().Fatalf("Something Wrong in File Reading Error:[%s]", err.Error())
		}
	}
	return &config
}

// GetString will return env value as string.
func (config *Config) GetString(key EnvKey) string {
	return config.viper.GetString(string(key))
}

// IsSet checks if environment variable is set.
func (config *Config) IsSet(key EnvKey) bool {
	return config.viper.IsSet(string(key))
}

// GetInt64 will return env value as int64
func (config *Config) GetInt64(key EnvKey) int64 {
	return config.viper.GetInt64(string(key))
}
