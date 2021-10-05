package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/litsea/logger"
	"github.com/lostsnow/cloudrain/telnet"
	"github.com/spf13/viper"
)

var (
	trace telnet.SessionTracer
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type wsWrapper struct {
	*websocket.Conn
}

func SetSessionTracer(t telnet.SessionTracer) {
	trace = t
}

func (wsw *wsWrapper) Write(p []byte) (n int, err error) {
	writer, err := wsw.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

func (wsw *wsWrapper) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := wsw.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func WebsocketHandler(c echo.Context) error {
	var err error
	var ip string
	ip = c.Request().Header.Get("X-REAL-IP")
	if ip == "" {
		ip, _, err = net.SplitHostPort(c.Request().RemoteAddr)
		if err != nil {
			return err
		}
	}

	up, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.Error("Error creating websocket: ", err)
		return err
	}

	logger.Info("Opening a proxy for ", c.Request().RemoteAddr)

	rw := io.ReadWriteCloser(&wsWrapper{up})
	onClose := func(s *telnet.Session) {
		if trace != nil {
			trace.Closed(s)
		}
	}
	up.SetCloseHandler(func(code int, text string) error {
		if err = up.Close(); err != nil {
			logger.Error(err)
		}
		return nil
	})

	var t telnet.Telnet
	err = viper.UnmarshalKey("telnet", &t)
	if err != nil {
		logger.Errorf("invalid telnet config: %v", t)
		return err
	}
	t.ClientIp = ip
	if t.Charset == "" {
		t.Charset = "utf-8"
	} else {
		t.Charset = strings.ToLower(t.Charset)
	}

	sess, err := t.NewSession(rw, onClose)
	if err == nil {
		if trace != nil {
			trace.Created(sess)
		}

		go handleCommand(up, sess)
	} else {
		logger.Errorf("error on session start: %s", err.Error())
		if err = up.Close(); err != nil {
			logger.Error(err)
		}
	}

	return nil
}

type command struct {
	Type    string
	Content string
}

// Send messages from the websocket to the telnet.
func handleCommand(ws *websocket.Conn, sess *telnet.Session) {
	for {
		_, bs, err := ws.ReadMessage()
		if err != nil {
			sess.Close()
			ws.Close()
			logger.Errorf("Error reading from ws(%s): %v", ws.RemoteAddr(), err)
			break
		}

		cmd := command{}
		if err = json.Unmarshal(bs, &cmd); err != nil {
			logger.Error(err)
			continue
		}

		t := cmd.Type
		switch t {
		case "cmd":
			sess.SendCommand(cmd.Content + "\n")

		case "naws":
			var w, h int
			_, err := fmt.Sscanf(cmd.Content, "%d,%d", &w, &h)
			if err == nil {
				sess.SendNaws(byte(w), byte(h))
			}

		case "atcp":
			sess.SendAtcp(cmd.Content)

		case "mxp":
			sess.SendMxp(cmd.Content)

		case "gmcp":
			sess.SendGmcp(strings.TrimSpace(cmd.Content))

		default:
			logger.Error(telnet.ErrInvalidCommand)
		}
	}
}
