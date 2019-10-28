package server

import (
	"strconv"

	"github.com/spf13/viper"
	"github.com/tehbilly/gmudc/telnet"
)

func newTelnet() (*telnet.Connection, error) {
	conn := telnet.New()

	addr := viper.GetString("telnet.host") + ":" + strconv.Itoa(viper.GetInt("telnet.port"))
	err := conn.Dial("tcp", addr)

	return conn, err
}
