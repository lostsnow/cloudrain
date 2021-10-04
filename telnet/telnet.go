package telnet

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/lostsnow/cloudrain/internal/version"
	"github.com/lostsnow/cloudrain/telnet/internal"
)

var (
	ErrClosed         = errors.New("channel closed")
	ErrInvalidCommand = errors.New("invalid command")
	ErrMaxConnection  = errors.New("maximum connection count reached")
)

var (
	TermTypes             = [][]byte{[]byte(version.AppName), []byte("MTTS 141")} // MTTS_ANSI | MTTS_UTF8 | MTTS_256_COLORS | MTTS_PROXY (1+4+8+128)
	MaxSessionConnections = 5
)

type Telnet struct {
	Debug        bool
	Host         string
	Port         int
	Charset      string
	GmcpSecret   string
	Secure       bool
	SecureVerify bool
	SendClientIp bool
	ClientIp     string
}

func (t *Telnet) Dial() (conn net.Conn, err error) {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	if t.Secure {
		config := tls.Config{
			InsecureSkipVerify: t.SecureVerify,
		}
		conn, err = tls.Dial("tcp", addr, &config)
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	return conn, err
}

func SetLogger(logger *log.Logger) {
	internal.Log = logger
}

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		internal.Log.Println(err)
	}
}
