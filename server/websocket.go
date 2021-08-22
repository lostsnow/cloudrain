package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/litsea/logger"
	"github.com/lostsnow/cloudrain/telnet"
	"github.com/spf13/viper"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	lock     sync.RWMutex
	sessions = make(map[string]*telnet.Session)
)

var (
	upgrader = websocket.Upgrader{}
)

type sessionTrace interface {
	SessionCreated()
	SessionClosed()
}

var trace sessionTrace

type wsWrapper struct {
	*websocket.Conn
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
	var sidCookie string
	cookie, err := c.Cookie("sessionid")
	if err == nil {
		sidCookie = cookie.Value
	}
	var tokenCookie string
	cookie, err = c.Cookie("token")
	if err == nil {
		tokenCookie = cookie.Value
	}

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
	me := telnet.NewMultiWriterEntry(rw)
	onClose := func(s *telnet.Session) {
		lock.Lock()
		defer lock.Unlock()
		delete(sessions, s.Id())
		logger.Infof("session ended %s.", plural(len(sessions)))
		if trace != nil {
			trace.SessionClosed()
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

	var sess *telnet.Session
	var sid string
	if t.MultiConnection {
		sid = c.QueryParam("sid")
		if sid == "" && sidCookie != "" {
			sid = sidCookie
		}

		if sid != "" && tokenCookie != "" {
			logger.Infof("try to attach session %s, %s.", sid, plural(len(sessions)))
			if attachToExistingSession(sid, tokenCookie, me) {
				go handleCommand(up, sessions[sid])
				return nil
			} else {
				sid = ""
			}
		}
	}

	sess, err = t.NewSession(sid, rw, onClose)
	if err == nil {
		sess.RemoteIp = ip
		lock.Lock()
		defer lock.Unlock()
		sessions[sess.Id()] = sess

		logger.Infof("session started %s, %s.", sess.Id(), plural(len(sessions)))

		if trace != nil {
			trace.SessionCreated()
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

func attachToExistingSession(sid, token string, me *telnet.MultiWriterEntry) bool {
	lock.RLock()
	defer lock.RUnlock()

	sess, ok := sessions[sid]
	if !ok {
		return false
	}

	if sess.Token() != token {
		logger.Errorf("invalid session %s token %s", sid, token)
		return false
	}

	err := sess.Attach(me)
	if err == nil {
		logger.Infof("session attached %s, %s", sid, plural(len(sessions)))
	} else {
		logger.Errorf("error on session attach %s: %s", sid, err.Error())
	}

	return true
}

func plural(value int) string {
	if value == 0 {
		return "(no active sessions)"
	} else if value == 1 {
		return "(1 active session)"
	} else {
		return "(" + strconv.Itoa(value) + " active sessions)"
	}
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
