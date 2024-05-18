package logger

import (
	"log"
	"os"
)

var InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
