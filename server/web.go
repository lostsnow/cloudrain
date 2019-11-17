package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func WebHandler(c *gin.Context) {
	sid := c.Request.URL.Query().Get("sid")
	u := url.URL{
		Scheme: viper.GetString("websocket.scheme"),
		Host:   fmt.Sprintf("%s:%d", viper.GetString("websocket.host"), viper.GetInt("websocket.port")),
		Path:   viper.GetString("websocket.path"),
	}

	if sid != "" {
		q := u.Query()
		q.Set("sid", sid)
		u.RawQuery = q.Encode()
	}

	c.JSON(
		http.StatusOK, gin.H{
			"websocketUrl": u.String(),
		},
	)
}
