package trello

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func VerifyTrelloSignature(callbackUrl, secret, headerHash string, body []byte) bool {
	content := append(body, []byte(callbackUrl)...)
	doubleHash := hmac256(content, secret)
	return doubleHash == headerHash
}

func hmac256(content []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
