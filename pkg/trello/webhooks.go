package trello

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// WebhookRequestBody represents the JSON structure of a Trello webhook request body
type WebhookRequestBody struct {
	Action struct {
		Type string `json:"type"`

		Display struct {
			TranslationKey string `json:"translationKey"`
		} `json:"display"`

		Data struct {
			Card struct {
				Id string `json:"id"`
			} `json:"card"`
		} `json:"data"`
	} `json:"action"`
}

// ParseArchivedCardId parses the archived card ID from the given webhook request body provided that
// the request body represents an archived card event. Returns empty string otherwise.
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

// VerifyWebhookSignature verifies the given Trello webhook signature (headerHash) by comparing it
// with a newly computed one using the webhook callback URL, Trello secret and the request body.
func VerifyWebhookSignature(callbackUrl, secret, headerHash string, body []byte) bool {
	content := append(body, []byte(callbackUrl)...)
	computedHash := hmacSha1(content, secret)
	return computedHash == headerHash
}

func hmacSha1(content []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
