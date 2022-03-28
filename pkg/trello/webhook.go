package trello

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/utkuufuk/entrello/internal/logger"
)

func VerifyTrelloSignature(callbackUrl, secret, headerHash string, body []byte) bool {
	content := append(body, []byte(callbackUrl)...)
	doubleHash := hmac256(content, secret)
	logger.Info("Header Hash:", headerHash)
	logger.Info("Double Hash:", doubleHash)
	return doubleHash == headerHash
}

func hmac256(content []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
