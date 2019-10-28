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
		if !viper.GetBool("debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.New()

		r.LoadHTMLGlob("templates/*.html")
		r.Static("/static", "./static")

		r.GET("/", server.WebHandler)
		r.GET("/"+viper.GetString("websocket.path"), server.WebsocketHandler)

		addr := viper.GetString("web.host") + ":" + strconv.Itoa(viper.GetInt("web.port"))
		log.Info("Listening and serving HTTP on ", addr)
		err := r.Run(addr)
		if err != nil {
			log.Fatal("Listening and serving HTTP error: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
