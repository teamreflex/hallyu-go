package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WebhookPayload struct {
	Content string `json:"content"`
}

func PostToDiscord(product RawProduct, client *http.Client) (bool, error) {
	// make request
	payload := WebhookPayload{
		Content: fmt.Sprintf("@here: %s: https://hallyusuperstore.com/products/%s", product.Title, product.Handle),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("DISCORD_WEBHOOK"), bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	// make sure request is good
	if err != nil {
		return false, fmt.Errorf("failed to create Discord request: %s", err)
	}

	// make request
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send Discord request: %s", err)
	}

	return res.StatusCode == 204, nil
}