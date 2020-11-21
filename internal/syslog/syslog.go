package syslog

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/utkuufuk/entrello/internal/config"
)

type Logger struct {
	enabled bool
	api     *tgbotapi.BotAPI
	chatId  int64
}

// NewLogger creates a new system logger instance
func NewLogger(cfg config.Telegram) (l Logger) {
	if !cfg.Enabled {
		return l
	}

	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Printf("could not create Telegram Logger: %v", err)
		return l
	}
	return Logger{true, api, cfg.ChatId}
}

// Debugf logs an debug message to stdout, but doesn't send a Telegram notification
func (l Logger) Debugf(msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	log.Println(msg)
}

// Printf logs an informational message to stdout, and also sends a Telegram notification if enabled
func (l Logger) Printf(msg string, v ...interface{}) {
	l.logf("Entrello:", msg, v...)
}

// Errorf logs an error message to stdout, and also sends a Telegram notification if enabled
func (l Logger) Errorf(msg string, v ...interface{}) {
	l.logf("Entrello Error:", msg, v...)
}

// Fatalf works like Errorf, but it returns with a non-zero exit code after logging
func (l Logger) Fatalf(msg string, v ...interface{}) {
	l.Errorf(msg, v...)
	os.Exit(1)
}

// logf prints the message to stdout, and after prepending the given prefix to the message,
// also sends a Telegram notification if enabled
func (l Logger) logf(prefix, msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	log.Println(msg)

	if !l.enabled || l.api == nil || l.chatId == 0 {
		return
	}
	msg = fmt.Sprintf("%s %s", prefix, msg)
	m := tgbotapi.NewMessage(l.chatId, msg)
	l.api.Send(m)
}
