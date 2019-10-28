package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func WebHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK, "index.html", gin.H{
			"websocketAddress": fmt.Sprintf("ws://%s:%d/%s",
				viper.GetString("websocket.host"), viper.GetInt("websocket.port"), viper.GetString("websocket.path")),
		},
	)
}
