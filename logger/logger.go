package logger

import (
	"log"
	"os"
)

var infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

func Info(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Error(format string, v ...any) {
	errorLog.Printf(format, v...)
}
