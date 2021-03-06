package cmd

import (
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/lostsnow/cloudrain/logger"
	"github.com/lostsnow/cloudrain/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve websocket and web frontend",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("web.debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.New()

		r.GET("/", server.WebHandler)
		r.GET("/"+viper.GetString("websocket.path"), server.WebsocketHandler)

		addr := viper.GetString("web.host") + ":" + strconv.Itoa(viper.GetInt("web.port"))
		log.Info("Listening on ", addr)
		err := r.Run(addr)
		if err != nil {
			log.Fatal("Listen error: ", err)
		}

		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic: %s", err)
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
