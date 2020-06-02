package syslog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/utkuufuk/entrello/internal/config"
)

func TestSystemLog(t *testing.T) {
	tt := []struct {
		name string
		cfg  config.Telegram
	}{
		{
			name: "null config",
			cfg:  config.Telegram{Enabled: false, Token: "", ChatId: 0},
		},
		{
			name: "psuedo-valid config but disabled",
			cfg: config.Telegram{
				Enabled: false,
				Token:   "xxxxxxxxxx:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				ChatId:  1234567890,
			},
		},
		{
			name: "psuedo-valid config and enabled",
			cfg: config.Telegram{
				Enabled: false,
				Token:   "xxxxxxxxxx:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				ChatId:  1234567890,
			},
		},
		{
			name: "enabled but invalid token",
			cfg: config.Telegram{
				Enabled: true,
				Token:   "banana",
				ChatId:  1234567890,
			},
		},
		{
			name: "enabled but invalid chat ID",
			cfg: config.Telegram{
				Enabled: true,
				Token:   "xxxxxxxxxx:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				ChatId:  0,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			msg := "test"
			logger := NewLogger(tc.cfg)
			var i, e bytes.Buffer

			log.SetOutput(&i)
			logger.Printf(msg)

			log.SetOutput(&e)
			logger.Errorf(msg)

			log.SetOutput(os.Stderr)

			oi := i.String()
			if !strings.Contains(oi, msg) {
				t.Errorf("wanted '%s' to contain '%s'", oi, msg)
			}

			oe := e.String()
			if !strings.Contains(oe, msg) {
				t.Errorf("wanted '%s' to contain '%s'", oe, msg)
			}
		})
	}
}
