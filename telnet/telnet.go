package telnet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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

type gmcpMessage struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value,omitempty"`
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

	pairs := bytes.Split(data, []byte{MSSP_VAR})
	msg := make(map[string]string)
	for _, pair := range pairs {
		ps := bytes.Split(pair, []byte{MSSP_VAL})
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
