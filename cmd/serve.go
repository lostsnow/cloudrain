package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/litsea/logger"
	"github.com/lostsnow/cloudrain/server"
	"github.com/lostsnow/cloudrain/telnet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve websocket and web frontend",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.Debug = viper.GetBool("web.debug")

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "/web/dist/",
			Filesystem: http.FS(assets.Web),
		}))

		server.SetSessionTracer(telnet.NewSessionTrace())

		e.GET("/app.config.js", configJsHandler)
		e.GET("/"+viper.GetString("websocket.path"), server.WebsocketHandler)

		addr := viper.GetString("web.host") + ":" + strconv.Itoa(viper.GetInt("web.port"))
		logger.Info("Listening on ", addr)
		e.Logger.Fatal(e.Start(addr))
	},
}

func configJsHandler(c echo.Context) error {
	js := fmt.Sprintf(`// eslint-disable-next-line no-unused-vars
const config = (() => {
  return {
    "VUE_APP_WEBSOCKET_URL": "%s://%s:%d/%s",
  };
})();`, viper.GetString("websocket.scheme"), viper.GetString("websocket.host"),
		viper.GetInt("websocket.port"), viper.GetString("websocket.path"))

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJavaScript)
	return c.HTML(http.StatusOK, js)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
