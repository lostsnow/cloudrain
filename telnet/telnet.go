package telnet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lostsnow/cloudrain/internal/version"
	"github.com/lostsnow/cloudrain/telnet/internal"
)

var (
	ErrClosed         = errors.New("channel closed")
	ErrInvalidCommand = errors.New("invalid command")
)

var (
	TermTypes = [][]byte{[]byte(version.AppName), []byte("MTTS 141")} // MTTS_ANSI | MTTS_UTF8 | MTTS_256_COLORS | MTTS_PROXY (1+4+8+128)
)

type Telnet struct {
	Debug        bool
	Host         string
	Port         int
	Charset      string
	GmcpSecret   string
	Secure       bool
	SendClientIp bool
	ClientIp     string
}

type message struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

type gmcpMessage struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value,omitempty"`
}

func (t *Telnet) Dial() (conn net.Conn, err error) {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	if t.Secure {
		conn, err = tls.Dial("tcp", addr, &tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	return conn, err
}

func SetLogger(logger *log.Logger) {
	internal.Log = logger
}

func (t *Telnet) NewSession(rwc io.ReadWriteCloser, onClose func(s *Session)) (sess *Session, err error) {
	me := NewWriterEntry(rwc)

	conn, err := t.Dial()
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 20)
	commandCh := make(chan interface{})
	sess = &Session{
		telnet:   t,
		conn:     conn,
		writer:   me,
		errCh:    errCh,
		buf:      make([]byte, 0, 1024),
		sbBuf:    make([]byte, 0, 64),
		debugBuf: make([]byte, 0, 256),
		commands: commandCh,
		once:     &sync.Once{},
		onClose:  onClose,
	}

	go func() {
		for {
			if sess.closed {
				break
			}
			if sess.gmcpHasHandShake {
				hello := map[string]string{
					"secret":  sess.telnet.GmcpSecret,
					"client":  version.AppName,
					"version": version.Version,
					"ip":      sess.telnet.ClientIp,
				}
				msg, _ := json.Marshal(hello)
				sess.SendGmcp("Core.Hello " + string(msg))
				break
			}
			time.Sleep(time.Millisecond * 50)
		}
	}()

	go func() {
		var err error

		defer sess.Close()
		quitCh := make(chan bool)
		defer close(quitCh)

		sess.initReadChannel(quitCh)

	ReadLoop:
		for {
			var timeout <-chan time.Time
			if len(sess.buf) > 0 {
				timeout = time.After(100 * time.Millisecond)
			} else {
				timeout = time.After(6 * time.Second)
			}

			select {
			case b, ok := <-sess.readCh:
				if !ok {
					break ReadLoop
				}
				err = sess.handleRead(b)
				break
			case <-timeout:
				err = sess.writeSocket()
				break
			case cmd := <-commandCh:
				err = sess.writeSocket()
				if err != nil {
					break
				}
				err = sess.handleCommand(cmd)
				break
			case err = <-errCh:
				break
			}

			if err != nil {
				break
			}
		}

		if err != nil {
			internal.Log.Println(err)
		} else {
			err = sess.writeSocket()
			if err != nil {
				internal.Log.Println(err)
			}
		}
	}()

	return sess, nil
}

func GMCPResponse(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	s := bytes.SplitN(data, []byte(" "), 2)

	var v interface{}
	if len(s) == 2 {
		err := json.Unmarshal(bytes.TrimSpace(s[1]), &v)
		if err != nil {
			internal.Log.Println(err)
			return nil, nil
		}
	}
	msg := gmcpMessage{
		Key:   string(s[0]),
		Value: v,
	}

	j, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func MSSPResponse(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	pairs := bytes.Split(data, []byte{MsspVar})
	msg := make(map[string]string)
	for _, pair := range pairs {
		ps := bytes.Split(pair, []byte{MsspVal})
		if len(ps) != 2 {
			continue
		}
		msg[string(ps[0])] = string(ps[1])
	}

	j, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return j, nil
}
