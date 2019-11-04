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
	OPT_TM          = 6   // timing mark. rfc860
	OPT_TTYPE       = 24  // terminal type. rfc930, rfc1091
	OPT_EOR         = 25  // end of record. rfc885
	OPT_NAWS        = 31  // negotiate about window size. rfc1073
	OPT_LINEMODE    = 34  // linemode. rfc1184
	OPT_ENVIRON     = 36  // environment option. rfc1408
	OPT_NEW_ENVIRON = 39  // new environment option. rfc1572
	OPT_CHARSET     = 42  // character set. rfc2066
	OPT_MSDP        = 69  // mud server data protocol. @see: https://tintin.sourceforge.io/protocols/msdp/
	OPT_MSSP        = 70  // mud server status protocol. @see: https://tintin.sourceforge.io/protocols/mssp/
	OPT_MCCP        = 86  // mud client compression protocol(v2). @see: https://tintin.sourceforge.io/protocols/mccp/
	OPT_MSP         = 90  // mud sound protocol. @see: https://www.zuggsoft.com/zmud/msp.htm
	OPT_MXP         = 91  // mud extension protocol. @see: https://www.zuggsoft.com/zmud/mxp.htm
	OPT_ATCP        = 200 // achaea telnet client protocol. @see: https://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html
	OPT_GMCP        = 201 // generic mud client protocol. @see: https://tintin.sourceforge.io/protocols/gmcp/

	// OPT_TTYPE
	TTYPE_IS   = 0
	TTYPE_SEND = 1

	// MTTS standard codes @see: https://tintin.sourceforge.io/protocols/mtts/
	MTTS_ANSI              = 1
	MTTS_VT100             = 2
	MTTS_UTF8              = 4
	MTTS_256_COLORS        = 8
	MTTS_MOUSE_TRACKING    = 16
	MTTS_OSC_COLOR_PALETTE = 32
	MTTS_SCREEN_READER     = 64
	MTTS_PROXY             = 128

	// OPT_ENVIRON, OPT_NEW_ENVIRON
	ENVIRON_IS      = 0
	ENVIRON_SEND    = 1
	ENVIRON_VAR     = 0
	ENVIRON_VALUE   = 1
	ENVIRON_ESC     = 2
	ENVIRON_USERVAR = 3

	// OPT_CHARSET
	CHARSET_REQUEST  = 1
	CHARSET_ACCEPTED = 2
	CHARSET_REJECTED = 3

	// OPT_MSSP
	MSDP_VAR         = 1
	MSDP_VAL         = 2
	MSDP_TABLE_OPEN  = 3
	MSDP_TABLE_CLOSE = 4
	MSDP_ARRAY_OPEN  = 5
	MSDP_ARRAY_CLOSE = 6

	// OPT_MSSP
	MSSP_VAR = 1
	MSSP_VAL = 2

	// OPT_LINEMODE
	LINEMODE_MODE          = 1
	LINEMODE_MODE_EDIT     = 1
	LINEMODE_MODE_TRAPSIG  = 2
	LINEMODE_MODE_ACK      = 4
	LINEMODE_MODE_SOFT_TAB = 8
	LINEMODE_MODE_LIT_ECHO = 16
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
