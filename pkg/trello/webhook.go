package trello

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

type WebhookRequestBody struct {
	Action Action `json:"action"`
}
type Action struct {
	Type    string `json:"type"`
	Display struct {
		TranslationKey string `json:"translationKey"`
	} `json:"display"`
	Data struct {
		Card struct {
			Id string `json:"id"`
		} `json:"card"`
	} `json:"data"`
}

func VerifyTrelloSignature(callbackUrl, secret, headerHash string, body []byte) bool {
	content := append(body, []byte(callbackUrl)...)
	doubleHash := hmac256(content, secret)
	return doubleHash == headerHash
}

func hmac256(content []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ParseArchivedCardId(body WebhookRequestBody) string {
	action := body.Action
	actionType := action.Type
	key := action.Display.TranslationKey
	id := action.Data.Card.Id
	if actionType == "updateCard" && key == "action_archived_card" {
		return id
	}
	return ""
}
