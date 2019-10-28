package server

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lostsnow/cloudrain/charset"
	log "github.com/lostsnow/cloudrain/logger"
	"github.com/spf13/viper"
	"github.com/tehbilly/gmudc/telnet"
)

type Server struct {
	w         http.ResponseWriter
	r         *http.Request
	Telnet    *telnet.Connection
	Websocket *websocket.Conn
	wg        sync.WaitGroup
	once      sync.Once
}

func TelnetProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/"+viper.GetString("websocket.path") {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	s := &Server{
		w: w,
		r: r,
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error creating websocket", 500)
		log.Error("Error creating websocket: ", err)
		return
	}
	defer c.Close()
	s.Websocket = c

	log.Info("Opening a proxy for ", r.RemoteAddr)
	t, err := newTelnet()
	if err != nil {
		log.Error("Error opening telnet proxy: ", err)
		return
	}
	defer t.Close()
	s.Telnet = t

	log.Infof("Connection open for %s. Proxying...", r.RemoteAddr)

	cs := strings.ToLower(viper.GetString("telnet.charset"))

	s.wg.Add(1)

	go s.writeMessage(cs)
	go s.readMessage(cs)

	// Wait until either go routine exits and then close both connections.
	s.wg.Wait()
	log.Info("Proxying completed for ", r.RemoteAddr)
}

// Send messages from the websocket to the telnet.
func (s *Server) writeMessage(cs string) {
	defer s.once.Do(func() { s.wg.Done() })
	for {
		_, bs, err := s.Websocket.ReadMessage()
		if err != nil {
			log.Errorf("Error reading from ws(%s): %v", s.r.RemoteAddr, err)
			break
		}

		if cs != "utf-8" {
			bs, err = charset.Encode(bs, cs)
			if err != nil {
				log.Error("Error convert websocket encoding")
				break
			}
		}

		// TODO: Partial writes.
		if _, err := s.Telnet.Write(bs); err != nil {
			log.Errorf("Error sending message to telnet for %s: %v", s.r.RemoteAddr, err)
			break
		}
	}
}

// Send messages from the telnet to the websocket.
func (s *Server) readMessage(cs string) {
	defer s.once.Do(func() { s.wg.Done() })
	br := bufio.NewReader(s.Telnet)
	for {
		bs := make([]byte, 1024)
		n, err := br.Read(bs)
		if err != nil {
			log.Errorf("Error reading from telnet for %s: %v", s.r.Host, err)
			break
		}

		bs = bytes.ReplaceAll(bs[:n], []byte{0xff, 0xf9}, []byte("\r"))
		if cs != "utf-8" {
			bs, err = charset.Decode(bs, cs)
			if err != nil {
				log.Error("Error convert telnet encoding")
				break
			}
		}

		if err = s.Websocket.WriteMessage(websocket.TextMessage, bs); err != nil {
			log.Errorf("Error sending to ws(%s): %v", s.r.RemoteAddr, err)
			break
		}
	}
}
