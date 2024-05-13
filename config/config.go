package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	configGlob config
	once       sync.Once
)

type config struct {
	PgUrl    string `mapstructure:"pg_url"`
	LogLevel string `mapstructure:"log_level"`
	HttpAddr string `mapstructure:"http_addr"`
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")       // path to look for the config file in
	viper.AddConfigPath("./../")    // path to look for the config file in
	viper.AddConfigPath("./../../") // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&configGlob)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
}

func Get() *config {
	once.Do(func() {
		initConfig()
	})

	return &configGlob
}
