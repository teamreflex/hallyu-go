package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	Content string `json:"content"`
}

func PostToDiscord(product RawProduct) (bool, error) {
	// build http client
	client := http.Client{
		Timeout: time.Second * 2,
	}

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
		str := fmt.Sprintf("Failed to create Discord request: %s", err)
		fmt.Println(str)
		return false, errors.New(str)
	}

	// make request
	res, err := client.Do(req)
	if err != nil {
		str:= fmt.Sprintf("Failed to send Discord request: %s", err)
		fmt.Println(str)
		return false, errors.New(str)
	}

	return res.StatusCode == 204, nil
}