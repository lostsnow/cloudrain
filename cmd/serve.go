package cmd

import (
	"github.com/spf13/cobra"
)

var websocketCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve websocket and web frontend",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(websocketCmd)
}
