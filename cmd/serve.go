package cmd

import (
	"net/http"
	"strconv"

	log "github.com/lostsnow/cloudrain/logger"
	"github.com/lostsnow/cloudrain/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve websocket and web frontend",
	Run: func(cmd *cobra.Command, args []string) {
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/", fs)
		http.HandleFunc("/"+viper.GetString("websocket.path"), server.TelnetProxy)

		addr := viper.GetString("web.host") + ":" + strconv.Itoa(viper.GetInt("web.port"))
		log.Info("Listen ", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal("Listen error: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
