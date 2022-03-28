package trello

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"

	"github.com/utkuufuk/entrello/internal/logger"
)

func VerifyTrelloSignature(callbackUrl, secret, headerHash string, body []byte) bool {
	content := append(body, []byte(callbackUrl)...)
	doubleHash := hmac256(content, secret)
	logger.Info("Header Hash: %s", headerHash)
	logger.Info("Double Hash: %s", doubleHash)
	return doubleHash == headerHash
}

func hmac256(content []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
