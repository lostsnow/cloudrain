package telnet

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/base32"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"github.com/lostsnow/cloudrain/charset"
	"github.com/lostsnow/cloudrain/telnet/internal"
)

const PingInterval = 5 * time.Second

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type setWindowSizeCommand struct{ width, height byte }
type atcpCommand string
type mxpCommand string

type Session struct {
	telnet     *Telnet
	conn       net.Conn
	id         string
	token      string
	writer     *multiWriter
	readCh     <-chan byte
	commands   chan<- interface{}
	errCh      chan<- error
	buf        []byte
	sbBuf      []byte
	latestBuf  []byte
	reader     *bufio.Reader
	termTypeIx int
	debugBuf   []byte
	lastWrite  time.Time
	RemoteIp   string
}

type message struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

type sessionIdentity struct {
	Sid   string `json:"sid"`
	Token string `json:"token"`
}

func (t *Telnet) NewSession(sid string, rwc io.ReadWriteCloser, onClose func(s *Session)) (sess *Session, err error) {
	if sid == "" {
		sid, err = generateId()
		if err != nil {
			return nil, err
		}
	}
	token, err := generateId()
	if err != nil {
		return nil, err
	}

	mw := NewMultiWriter(rwc)

	conn, err := t.Dial()
	if err != nil {

		return nil, err
	}

	errCh := make(chan error, 20)
	commandCh := make(chan interface{})
	sess = &Session{
		telnet:   t,
		conn:     conn,
		writer:   mw,
		errCh:    errCh,
		buf:      make([]byte, 0, 1024),
		sbBuf:    make([]byte, 0, 64),
		reader:   bufio.NewReader(conn),
		debugBuf: make([]byte, 0, 256),
		id:       sid,
		token:    token,
		commands: commandCh,
	}

	si := &sessionIdentity{
		Sid:   sid,
		Token: token,
	}
	siBytes, err := json.Marshal(si)
	if err != nil {
		return nil, err
	}

	go func() {
		var err error

		defer Close(conn)
		defer mw.Close()
		defer close(commandCh)

		quitCh := make(chan bool)
		defer close(quitCh)

		sess.initReadChannel(quitCh)
		if err := sess.writeSocketRaw("session", siBytes); err != nil {
			internal.Log.Println(err)
		}

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

		if onClose != nil {
			onClose(sess)
		}
	}()

	return sess, nil
}

func (sess *Session) Attach(me *MultiWriterEntry) error {
	if err := sess.writer.attach(me.writer); err != nil {
		if e := me.writer.Close(); e != nil {
			internal.Log.Println(e)
		}
		return err
	}

	if len(sess.latestBuf) > 0 {
		msg := &message{
			Event:   "text",
			Content: string(sess.latestBuf),
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, err = sess.writer.Write(msgBytes)
		if err != nil {
			return err
		}
		_, err = sess.writer.Write([]byte(""))
		if err != nil {
			return err
		}
	}

	return nil
}

func (sess *Session) handleRead(b byte) error {
	if b == IAC {
		return sess.handleOption()
	} else {
		sess.buf = append(sess.buf, b)
	}
	return nil
}

func (sess *Session) readOptionByte(printable bool) (byte, error) {
	b, ok := <-sess.readCh
	if !ok {
		return 0, ErrClosed
	}
	sess.DumpOptionByte(b, printable)
	return b, nil
}

func (sess *Session) handleOption() error {
	var err error
	sess.DumpOptionByte(IAC, false)
	first, err := sess.readOptionByte(false)
	if err != nil {
		return err
	}
	switch first {
	case EOR:
		err = sess.writeSocket()
		if err != nil {
			return err
		}
		sess.FlushTelnetBuffer()

	case SB:
		buf := sess.sbBuf[0:0]
		second, err := sess.readOptionByte(false)
		if err != nil {
			return err
		}
		for {
			if second == IAC {
				second, err := sess.readOptionByte(true)
				if err != nil {
					return err
				}
				if second == SE {
					break
				}
				buf = append(buf, IAC)
			} else {
				buf = append(buf, second)
				second, err = sess.readOptionByte(true)
				if err != nil {
					return err
				}
			}
		}
		sess.FlushTelnetBuffer()
		return sess.handleSb(buf)

	case DO:
		second, err := sess.readOptionByte(false)
		if err != nil {
			return err
		}
		sess.FlushTelnetBuffer()
		return sess.handleDo(second)

	case DONT:
		_, err = sess.readOptionByte(false)
		if err != nil {
			return err
		}
		sess.FlushTelnetBuffer()

	case WILL:
		second, err := sess.readOptionByte(false)
		if err != nil {
			return err
		}
		sess.FlushTelnetBuffer()
		switch second {
		case OPT_MSSP:
			return sess.handleDo(second)
		default:
			return sess.handleWill(second)
		}

	case WONT:
		_, err = sess.readOptionByte(false)
		if err != nil {
			return err
		}
		sess.FlushTelnetBuffer()

	default:
		sess.FlushTelnetBuffer()
	}

	return nil
}

func (sess *Session) handleSb(data []byte) error {
	var err error
	if len(data) < 1 {
		return nil
	}

	option := data[0]

	switch option {
	case OPT_TTYPE:
		if len(data) != 2 || data[1] != TTYPE_SEND {
			return nil
		}
		err = sess.writeSb(OPT_TTYPE, []byte{TTYPE_IS}, TermTypes[sess.termTypeIx])
		if err != nil {
			return err
		}
		sess.termTypeIx = (sess.termTypeIx + 1) % len(TermTypes)

	case OPT_NEW_ENVIRON, OPT_ENVIRON:
		if len(data) < 2 || data[1] != ENVIRON_SEND {
			return nil
		}
		return sess.writeSb(option, []byte{ENVIRON_IS})

	case OPT_LINEMODE:
		if len(data) != 3 || data[1] != LINEMODE_MODE {
			return nil
		}
		mask := data[2]
		if mask&LINEMODE_MODE_ACK == LINEMODE_MODE_ACK {
			return nil
		}

		replyMask := mask | LINEMODE_MODE_ACK
		return sess.writeSb(OPT_LINEMODE, []byte{LINEMODE_MODE, replyMask})

	case OPT_MSSP:
		return sess.writeSocketRaw("mssp", data[1:])

	case OPT_ATCP:
		if string(data[1:]) == "Auth.Request ON" && sess.telnet.SendRemoteIp {
			sess.sendATCPRemoteIp()
		}
		return sess.writeSocketRaw("atcp", data[1:])

	case OPT_MXP:
		return sess.writeSocketRaw("mxp", data[1:])
	}

	return nil
}

func (sess *Session) writeSb(option byte, data ...[]byte) error {
	size := 0
	for _, d := range data {
		size += len(d)
	}

	buf := make([]byte, 0, size+5)

	buf = append(buf, IAC, SB, option)
	sess.TelnetDebug("<- ")
	sess.DumpOptionByte(IAC, false)
	sess.DumpOptionByte(SB, false)
	sess.DumpOptionByte(option, false)

	for _, d := range data {
		buf = append(buf, d...)
		for _, b := range d {
			sess.DumpOptionByte(b, true)
		}
	}

	buf = append(buf, IAC, SE)
	sess.DumpOptionByte(IAC, false)
	sess.DumpOptionByte(SE, false)
	sess.FlushTelnetBuffer()

	_, err := sess.conn.Write(buf)
	return err
}

func (sess *Session) writeSbString(option byte, text string) error {
	return sess.writeSb(option, []byte(text))
}

func (sess *Session) sendATCPRemoteIp() {
	addr := sess.RemoteIp
	portSep := strings.Index(addr, ":")
	if portSep != -1 {
		addr = addr[:portSep]
	}
	names, err := net.LookupAddr(addr)

	var name string
	if err == nil && len(names) > 0 {
		name = names[0]
	} else {
		name = addr
	}
	name = strings.TrimSuffix(name, ".")
	err = sess.writeSbString(OPT_ATCP, "ava_remoteip "+addr+" "+name)
}

func (sess *Session) handleDo(second byte) error {
	var err error
	switch second {
	case OPT_ATCP:
		err = sess.writeOption(IAC, WILL, second)
		if err != nil {
			return err
		}

		var options = []string{
			"hello CloudRain 0.1.0 ALPHA",
			"ping 1",
			"keepalive 1",
			"composer 1",
			"topvote 1",
			"auth 0",
			"char_name 1",
			"char_vitals 1",
			"room_brief 1",
			"room_exits 1",
			"map_display 1",
		}

		err = sess.writeSbString(OPT_ATCP, strings.Join(options, "\n"))
		if err != nil {
			return err
		}

	case OPT_NAWS:
		return sess.writeOption(IAC, WILL, second)

	case OPT_TTYPE, OPT_ENVIRON, OPT_LINEMODE, OPT_NEW_ENVIRON:
		return sess.writeOption(IAC, WILL, second)

	case OPT_TM:
		err = sess.writeSocket()
		if err != nil {
			return err
		}
		return sess.writeOption(IAC, WILL, second)

	case OPT_MSSP:
		return sess.writeOption(IAC, DO, second)

	default:
		return sess.writeOption(IAC, WONT, second)
	}

	return nil
}

func (sess *Session) handleWill(option byte) error {
	return sess.writeOption(IAC, DONT, option)
}

func (sess *Session) handleCommand(cmd interface{}) error {
	switch cmd := cmd.(type) {
	case string:
		data := []byte(cmd)
		data, err := charset.Encode(data, sess.telnet.Charset)
		if err != nil {
			return nil
		}

		_, err = sess.conn.Write(data)
		return err

	case setWindowSizeCommand:
		return sess.writeOption(IAC, SB, OPT_NAWS, 0, cmd.width, 0, cmd.height, IAC, SE)

	case atcpCommand:
		return sess.writeSb(OPT_ATCP, []byte(cmd))

	case mxpCommand:
		return sess.writeSb(OPT_MXP, []byte(cmd))

	default:
		return ErrInvalidCommand
	}
}

func (sess *Session) SendCommand(command string) {
	sess.commands <- command
}

func (sess *Session) SendNaws(width, height byte) {
	sess.commands <- setWindowSizeCommand{width, height}
}

func (sess *Session) SendAtcp(text string) {
	sess.commands <- atcpCommand(text)
}

func (sess *Session) SendMxp(text string) {
	sess.commands <- mxpCommand(text)
}

func (sess *Session) initReadChannel(quitCh <-chan bool) {
	ch := make(chan byte)
	compressionSequence := []byte{IAC, SB, OPT_MCCP, IAC, SE}
	var seqix int
	compressionStarted := false
	reader := sess.reader
	go func() {
		defer close(ch)
		for {
			b, err := reader.ReadByte()

			if err == io.EOF {
				return
			} else if err != nil {
				sess.errCh <- err
				return
			}

			select {
			case ch <- b:
				break
			case <-quitCh:
				return
			}

			if !compressionStarted && b == compressionSequence[seqix] {
				seqix++
				if seqix == len(compressionSequence) {
					seqix = 0
					zReader, err := zlib.NewReader(sess.reader)
					if err != nil {
						internal.Log.Println(err)
						if e := sess.writeOption(IAC, DONT, OPT_MCCP); e != nil {
							internal.Log.Println(err)
						}
					} else {
						reader = bufio.NewReader(zReader)
						compressionStarted = true
					}
				}
			} else {
				seqix = 0
			}
		}
	}()
	sess.readCh = ch
}

func (sess *Session) writeSocket() error {
	err := sess.writeSocketRaw("text", sess.buf)
	if err != nil {
		sess.errCh <- err
		return err
	}
	sess.buf = sess.buf[0:0]
	return nil
}

func (sess *Session) writeSocketRaw(event string, data []byte) error {
	if len(data) == 0 {
		return sess.testPing()
	}
	sess.lastWrite = time.Now()

	writer := sess.writer

	data, err := charset.Decode(data, sess.telnet.Charset)
	if err != nil {
		return err
	}

	if event == "text" {
		keepLines := 1000
		latestBuf := append(sess.latestBuf, data...)
		bs := bytes.Split(latestBuf, []byte("\n"))
		if len(bs) <= keepLines {
			sess.latestBuf = latestBuf
		} else {
			sess.latestBuf = bytes.Join(bs[len(bs)-keepLines:], []byte("\n"))
		}
	}

	msg := &message{
		Event:   event,
		Content: string(data),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = writer.Write(msgBytes)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(""))
	if err != nil {
		return err
	}

	return nil
}

func (sess *Session) writeOption(data ...byte) error {
	_, err := sess.conn.Write(data)
	if err != nil {
		return err
	}
	sess.TelnetDebug("<- ")
	for _, b := range data {
		sess.DumpOptionByte(b, false)
	}
	sess.FlushTelnetBuffer()

	return nil
}

func (sess *Session) TelnetDebug(text ...string) {
	if !sess.telnet.Debug {
		return
	}
	for _, txt := range text {
		sess.debugBuf = append(sess.debugBuf, []byte(txt)...)
	}
}

func (sess *Session) DumpOptionByte(b byte, printable bool) {
	if b >= 240 {
		sess.TelnetDebug("<", CmdNames[b], ">")
	} else if printable && b >= 32 && b <= 126 {
		sess.TelnetDebug(string([]byte{b}))
	} else {
		sess.TelnetDebug("<", strconv.FormatInt(int64(b), 10), ">")
	}
}

func (sess *Session) FlushTelnetBuffer() {
	if !sess.telnet.Debug {
		return
	}
	if len(sess.debugBuf) > 0 {
		internal.TelnetDebug.Println(string(sess.debugBuf))
		sess.debugBuf = sess.debugBuf[0:0]
	}
}

func (sess *Session) testPing() error {
	now := time.Now()
	if now.After(sess.lastWrite.Add(PingInterval)) {
		msg := &message{
			Event:   "ping",
			Content: "",
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		_, err = sess.writer.Write(msgBytes)
		if err != nil {
			return err
		}

		sess.lastWrite = now
	}
	return nil
}

func (sess *Session) Token() string {
	return sess.token
}

func (sess *Session) Id() string {
	return sess.id
}

func generateId() (string, error) {
	idEncoding := base32.NewEncoding("123456789abcdefghijklmnopqrstuvw")
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil
	}
	return strings.TrimRight(idEncoding.EncodeToString(b), "="), nil
}
