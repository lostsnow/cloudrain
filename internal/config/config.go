package config

import (
	"fmt"
	"os"

	"github.com/litsea/logger"
	"github.com/spf13/viper"
)

func ReadConfig(cfgFile, configPath string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath(configPath)
	}

	return viper.ReadInConfig()
}

func InitLogger() {
	v := viper.Sub("log")

	driver := viper.GetString("log-driver")
	if err := logger.NewLogger(v, driver); err != nil {
		fmt.Println("init logger failed: ", err)
		os.Exit(1)
	}

	logger.Info("Logger initialized")
}
