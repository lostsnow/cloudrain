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
)

func TelnetProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/"+viper.GetString("websocket.path") {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
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

	log.Info("Opening a proxy for ", r.RemoteAddr)
	t, err := newTelnet()
	if err != nil {
		log.Error("Error opening telnet proxy: ", err)
		return
	}
	defer t.Close()

	log.Infof("Connection open for %s. Proxying...", r.RemoteAddr)

	cs := strings.ToLower(viper.GetString("telnet.charset"))
	var wg sync.WaitGroup
	var once sync.Once
	wg.Add(1)

	// Send messages from the websocket to the telnet.
	go func() {
		defer once.Do(func() { wg.Done() })
		for {
			_, bs, err := c.ReadMessage()
			if err != nil {
				log.Errorf("Error reading from ws(%s): %v", r.RemoteAddr, err)
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
			if _, err := t.Write(bs); err != nil {
				log.Errorf("Error sending message to telnet for %s: %v", r.RemoteAddr, err)
				break
			}
		}
	}()

	// Send messages from the telnet to the websocket.
	go func() {
		defer once.Do(func() { wg.Done() })
		br := bufio.NewReader(t)
		for {
			bs := make([]byte, 1024)
			n, err := br.Read(bs)
			if err != nil {
				log.Errorf("Error reading from telnet for %s: %v", r.Host, err)
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

			if err = c.WriteMessage(websocket.TextMessage, bs); err != nil {
				log.Errorf("Error sending to ws(%s): %v", r.RemoteAddr, err)
				break
			}
		}
	}()

	// Wait until either go routine exits and then close both connections.
	wg.Wait()
	log.Info("Proxying completed for ", r.RemoteAddr)
}
