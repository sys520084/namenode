package config

import (
	"log"

	"github.com/spf13/viper"
)

func Init(configName string) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	viper.SetDefault("debug", true)
	viper.SetDefault("log_json_format", false)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("Fatal error config file: %s \n", err)
	}

}
