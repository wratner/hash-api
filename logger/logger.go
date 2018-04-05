package logger

import (
	"io"
	"log"
)

// Global logging variables for info and error level
var (
	Info  *log.Logger
	Error *log.Logger
)

// Init initializes the loggers with specific formats.
// These format is "LEVEL: " DATE | TIME | FILE AND LINE NUMBER
func Init(infoHandle io.Writer, errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
