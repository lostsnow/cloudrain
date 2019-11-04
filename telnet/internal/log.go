package internal

import (
	"log"
	"os"
)

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

var Log = Logger(log.New(os.Stderr, "[telnet] ", log.Ldate|log.Ltime|log.Lshortfile))
var TelnetDebug = Logger(log.New(os.Stdout, "", 0))
