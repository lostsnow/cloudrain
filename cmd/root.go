package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cloudrain",
	Short: "CloudRain",
	Long:  `CloudRain`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
	}

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %v\n", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
		os.Exit(-1)
	}
}
