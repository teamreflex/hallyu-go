package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RawProduct struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Handle string `json:"handle"`
	CreatedAt string `json:"created_at"`
}

type RawProductResponse struct {
	Products []RawProduct `json:"products"`
}

var url = "https://hallyusuperstore.com/products.json"

func GetProducts(client *http.Client) ([]RawProduct, error) {
	// make request
	req, err := http.NewRequest(http.MethodGet, url, nil)

	// make sure request is good
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// make request
	fmt.Println("Fetching products...")
	res, err := client.Do(req)

	// ensure response is good
	if (res.StatusCode != 200) {
		if res.StatusCode == 429 {
			fmt.Println("Rate limit hit, waiting 5 minutes...")
			time.Sleep(time.Minute * 5)
			return GetProducts(client)
		}

		return nil, fmt.Errorf("failed to get products: %s", err)
	}

	// read response body
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %s", err)
	}

	// parse response
	response := RawProductResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %s", jsonErr)
	}

	return response.Products, nil
}