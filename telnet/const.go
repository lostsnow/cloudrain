package telnet

const (
	// Telnet commands. rfc854, rfc1116, rfc1123
	EOF   = 236 // end of file
	SUSP  = 237 // suspend process
	ABORT = 238 // abort process
	EOR   = 239 // end of record (transparent mode, used for prompt marking)
	SE    = 240 // end sub negotiation
	NOP   = 241 // nop (used for keep alive messages	)
	DM    = 242 // data mark--for connect. cleaning
	BREAK = 243 // break
	IP    = 244 // interrupt process (permanently)
	AO    = 245 // abort output (but let program finish)
	AYT   = 246 // are you there
	EC    = 247 // erase the current character
	EL    = 248 // erase the current line
	GA    = 249 // you may reverse the line (used for prompt marking)
	SB    = 250 // interpret as subnegotiation
	WILL  = 251 // I will use option
	WONT  = 252 // I won"t use option
	DO    = 253 // please, you use option
	DONT  = 254 // you are not to use option
	IAC   = 255 // interpret as command

	// Telnet Options. rfc855
	OptTM         = 6   // timing mark. rfc860
	OptTType      = 24  // terminal type. rfc930, rfc1091
	OptEOR        = 25  // end of record. rfc885
	OptNAWS       = 31  // negotiate about window size. rfc1073
	OptLineMode   = 34  // linemode. rfc1184
	OptEnviron    = 36  // environment option. rfc1408
	OptNewEnviron = 39  // new environment option. rfc1572
	OptCharset    = 42  // character set. rfc2066
	OptMSDP       = 69  // mud server data protocol. @see: https://tintin.sourceforge.io/protocols/msdp/
	OptMSSP       = 70  // mud server status protocol. @see: https://tintin.sourceforge.io/protocols/mssp/
	OptMCCP       = 86  // mud client compression protocol(v2). @see: https://tintin.sourceforge.io/protocols/mccp/
	OptMSP        = 90  // mud sound protocol. @see: https://www.zuggsoft.com/zmud/msp.htm
	OptMXP        = 91  // mud extension protocol. @see: https://www.zuggsoft.com/zmud/mxp.htm
	OptATCP       = 200 // achaea telnet client protocol. @see: https://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html
	OptGMCP       = 201 // generic mud client protocol. @see: https://tintin.sourceforge.io/protocols/gmcp/

	// OptTType
	TTypeIs   = 0
	TTypeSend = 1

	// MTTS standard codes @see: https://tintin.sourceforge.io/protocols/mtts/
	TTypeANSI            = 1
	TTypeVT100           = 2
	TTypeUTF8            = 4
	TType256Colors       = 8
	TTypeMouseTracking   = 16
	TTypeOscColorPalette = 32
	TTypeScreenReader    = 64
	TTypeProxy           = 128
	TTypeTrueColor       = 256
	TTypeMNES            = 512
	TTypeMSLP            = 1024

	// OptEnviron, OptNewEnviron
	EnvironIs      = 0
	EnvironSend    = 1
	EnvironVar     = 0
	EnvironValue   = 1
	EnvironESC     = 2
	EnvironUserVar = 3

	// OptCharset
	CharsetRequest  = 1
	CharsetAccepted = 2
	CharsetRejected = 3

	// OptMSDP
	MSDPVar        = 1
	MSDPVal        = 2
	MSDPTableOpen  = 3
	MSDPTableClose = 4
	MSDPArrayOpen  = 5
	MSDPArrayClose = 6

	// OPT_MSSP
	MsspVar = 1
	MsspVal = 2

	// OptLineMode
	LineMode        = 1
	LineModeEdit    = 1
	LineModeTrapSig = 2
	LineModeAck     = 4
	LineModeSoftTab = 8
	LineModeLitEcho = 16
)

var CmdNames = map[byte]string{
	EOF:   "EOF",
	SUSP:  "SUSP",
	ABORT: "ABORT",
	EOR:   "EOR",
	SE:    "SE",
	NOP:   "NOP",
	DM:    "DM",
	BREAK: "BREAK",
	IP:    "IP",
	AO:    "AO",
	AYT:   "AYT",
	EC:    "EC",
	EL:    "EL",
	GA:    "GA",
	SB:    "SB",
	WILL:  "WILL",
	WONT:  "WONT",
	DO:    "DO",
	DONT:  "DONT",
	IAC:   "IAC",
}
