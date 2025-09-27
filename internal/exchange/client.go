// internal/exchange/client.go
package exchange

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.bitunix.com"
)

type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) PlaceOrder(ctx context.Context, symbol string, side string, orderType string, qty float64) (string, error) {
	endpoint := "/api/v1/trade/order"
	method := "POST"

	body := map[string]interface{}{
		"symbol":   symbol,
		"side":     side,
		"type":     orderType,
		"quantity": qty,
	}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	signature := c.createSignature(timestamp, method, endpoint, string(bodyBytes))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("B-API-KEY", c.apiKey)
	req.Header.Set("B-TIMESTAMP", timestamp)
	req.Header.Set("B-SIGN", signature)

	log.Printf("Sending order to Bitunix: %s", string(bodyBytes))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Successfully placed order. Response: %s", string(respBody))
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	orderID := fmt.Sprintf("%v", result["data"])

	return orderID, nil
}

func (c *Client) createSignature(timestamp, method, endpoint, body string) string {
	preHashString := timestamp + c.apiKey + method + endpoint + body

	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(preHashString))
	return hex.EncodeToString(h.Sum(nil))
}
