package telnet

import (
	"bufio"
	"compress/zlib"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/lostsnow/cloudrain/charset"
	"github.com/lostsnow/cloudrain/telnet/internal"
)

type setWindowSizeCommand struct{ width, height byte }
type atcpCommand string
type mxpCommand string
type gmcpCommand string

type Session struct {
	telnet           *Telnet
	conn             net.Conn
	writer           *writerEntry
	readCh           <-chan byte
	commands         chan<- interface{}
	errCh            chan<- error
	buf              []byte
	sbBuf            []byte
	termTypeIx       int
	debugBuf         []byte
	once             *sync.Once
	closed           bool
	onClose          func(s *Session)
	gmcpHasHandShake bool
}

func (sess *Session) Close() {
	sess.once.Do(func() {
		sess.closed = true
		close(sess.commands)
		sess.writer.Close()
		_ = sess.conn.Close()

		if sess.onClose != nil {
			sess.onClose(sess)
		}
	})
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
		case OptMSSP:
			return sess.handleDo(second)
		case OptGMCP:
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
	case OptTType:
		if len(data) != 2 || data[1] != TTypeSend {
			return nil
		}
		err = sess.writeSb(OptTType, []byte{TTypeIs}, TermTypes[sess.termTypeIx])
		if err != nil {
			return err
		}
		sess.termTypeIx = (sess.termTypeIx + 1) % len(TermTypes)
	case OptNewEnviron, OptEnviron:
		if len(data) < 2 || data[1] != EnvironSend {
			return nil
		}
		return sess.writeSb(option, []byte{EnvironIs},
			[]byte{EnvironVar}, []byte("REAL_IP"), []byte{EnvironValue}, []byte(sess.telnet.ClientIp),
		)
	case OptLineMode:
		if len(data) != 3 || data[1] != LineMode {
			return nil
		}
		mask := data[2]
		if mask&LineModeAck == LineModeAck {
			return nil
		}

		replyMask := mask | LineModeAck
		return sess.writeSb(OptLineMode, []byte{LineMode, replyMask})
	case OptMSSP:
		md, err := MSSPResponse(data[1:])
		if err != nil {
			internal.Log.Println(err)
			return nil
		}
		return sess.writeSocketRaw("mssp", md)
	case OptATCP:
		if string(data[1:]) == "Auth.Request ON" && sess.telnet.SendClientIp {
			sess.sendATCPClientIp()
		}
		return sess.writeSocketRaw("atcp", data[1:])
	case OptGMCP:
		gd, err := GMCPResponse(data[1:])
		if err != nil {
			internal.Log.Println(err)
			return nil
		}
		if err := sess.writeSocketRaw("gmcp", gd); err != nil {
			internal.Log.Println(err)
		}
	case OptMXP:
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

func (sess *Session) sendATCPClientIp() {
	addr := sess.telnet.ClientIp
	names, err := net.LookupAddr(addr)

	var name string
	if err == nil && len(names) > 0 {
		name = names[0]
	} else {
		name = addr
	}
	name = strings.TrimSuffix(name, ".")
	err = sess.writeSbString(OptATCP, "ava_remoteip "+addr+" "+name)
	if err != nil {
		internal.Log.Println(err)
	}
}

func (sess *Session) handleDo(second byte) error {
	var err error
	switch second {
	case OptATCP:
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

		err = sess.writeSbString(OptATCP, strings.Join(options, "\n"))
		if err != nil {
			return err
		}
	case OptNAWS:
		return sess.writeOption(IAC, WILL, second)
	case OptTType, OptEnviron, OptLineMode, OptNewEnviron:
		return sess.writeOption(IAC, WILL, second)
	case OptTM:
		err = sess.writeSocket()
		if err != nil {
			return err
		}
		return sess.writeOption(IAC, WILL, second)
	case OptMSSP:
		return sess.writeOption(IAC, DO, second)
	case OptGMCP:
		err = sess.writeOption(IAC, DO, second)
		if err != nil {
			return err
		}
		sess.gmcpHasHandShake = true
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
		if sess.telnet.Charset != "utf-8" {
			d, err := charset.Encode(data, sess.telnet.Charset)
			if err != nil {
				return nil //nolint:nilerr
			}
			data = d
		}

		_, err := sess.conn.Write(data)
		return err
	case setWindowSizeCommand:
		return sess.writeOption(IAC, SB, OptNAWS, 0, cmd.width, 0, cmd.height, IAC, SE)
	case atcpCommand:
		return sess.writeSb(OptATCP, []byte(cmd))
	case mxpCommand:
		return sess.writeSb(OptMXP, []byte(cmd))
	case gmcpCommand:
		return sess.writeSb(OptGMCP, []byte(cmd))
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

func (sess *Session) SendGmcp(text string) {
	sess.commands <- gmcpCommand(text)
}

func (sess *Session) initReadChannel(quitCh <-chan bool) {
	ch := make(chan byte)
	compressionSequence := []byte{IAC, SB, OptMCCP, IAC, SE}
	var seqix int
	compressionStarted := false
	reader := bufio.NewReader(sess.conn)
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

			if compressionStarted || b != compressionSequence[seqix] {
				seqix = 0
				continue
			}

			seqix++
			if seqix != len(compressionSequence) {
				continue
			}
			seqix = 0
			zReader, err := zlib.NewReader(reader)
			if err == nil {
				reader = bufio.NewReader(zReader)
				compressionStarted = true
				continue
			}
			internal.Log.Println(err)
			if e := sess.writeOption(IAC, DONT, OptMCCP); e != nil {
				internal.Log.Println(err)
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
		return nil
	}

	writer := sess.writer

	if sess.telnet.Charset != "utf-8" {
		d, err := charset.Decode(data, sess.telnet.Charset)
		if err != nil {
			return err
		}
		data = d
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
