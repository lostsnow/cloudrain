package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/lostsnow/cloudrain/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	assets  *Assets
)

var rootCmd = &cobra.Command{
	Use:   "cloudrain",
	Short: "CloudRain",
	Long:  `CloudRain`,
}

type Assets struct {
	Web embed.FS
}

func Execute(as *Assets) {
	assets = as
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/app.yml)")
}

func initConfig() {
	if err := config.ReadConfig(cfgFile, "./configs"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Using config file: %v\n", viper.ConfigFileUsed())
	config.InitLogger()
}
