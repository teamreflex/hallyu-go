package main

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

func GetProducts() []RawProduct {
	// build http client
	client := http.Client{
		Timeout: time.Second * 2,
	}

	// make request
	req, err := http.NewRequest(http.MethodGet, url, nil)

	// make sure request is good
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err)
		return GetProducts()
	}

	// make request
	fmt.Println("Fetching products...")
	res, err := client.Do(req)

	// ensure response is good
	if (res.StatusCode != 200) {
		if res.StatusCode == 429 {
			fmt.Println("Rate limit hit, waiting 5 minutes...")
			time.Sleep(time.Minute * 5)
			return GetProducts()
		}

		fmt.Printf("Failed to get products: %s\n", err)
		time.Sleep(time.Minute * 1)
		return GetProducts()
	}

	// read response body
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to parse response body: %s\n", err)
		time.Sleep(time.Minute * 1)
		return GetProducts()
	}

	// parse response
	response := RawProductResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		fmt.Printf("Failed to parse response JSON: %s\n", jsonErr)
		time.Sleep(time.Second * 10)
		return GetProducts()
	}

	return response.Products
}