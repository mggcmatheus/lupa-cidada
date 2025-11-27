package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient é um cliente HTTP com rate limiting
type HTTPClient struct {
	client      *http.Client
	rateLimiter <-chan time.Time
}

// NewHTTPClient cria um novo cliente HTTP com rate limiting
func NewHTTPClient(requestsPerSecond int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: time.Tick(time.Second / time.Duration(requestsPerSecond)),
	}
}

// Get faz uma requisição GET e decodifica o JSON
func (c *HTTPClient) Get(url string, result interface{}) error {
	<-c.rateLimiter // Rate limiting

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "LupaCidada/1.0 (Portal de Transparência)")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return nil
}
