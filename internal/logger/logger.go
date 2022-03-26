package logger

import (
	"fmt"
	"log"
)

func Info(msg string, v ...interface{}) {
	logf("[INFO]", msg, v...)
}

func Warn(msg string, v ...interface{}) {
	logf("[WARN]", msg, v...)
}

func Error(msg string, v ...interface{}) {
	logf("[ERROR]", msg, v...)
}

func logf(prefix, msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	log.Printf("%s %s\n", prefix, msg)
}
