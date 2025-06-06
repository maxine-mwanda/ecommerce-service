package notifications

import (
	"context"
	"ecommerce-service/internal/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SMSSender interface {
	Send(ctx context.Context, phone, message string) error
}

type AfricaTalkingSMSSender struct {
	cfg config.SMSConfig
}

func NewAfricaTalkingSMSSender(cfg config.SMSConfig) *AfricaTalkingSMSSender {
	return &AfricaTalkingSMSSender{cfg: cfg}
}

func (s *AfricaTalkingSMSSender) Send(ctx context.Context, phone, message string) error {
	apiURL := "https://api.africastalking.com/version1/messaging"

	data := url.Values{}
	data.Set("username", s.cfg.Username)
	data.Set("to", phone)
	data.Set("message", message)
	if s.cfg.SenderID != "" {
		data.Set("from", s.cfg.SenderID)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", s.cfg.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("AfricaTalking SMS response: %s\n", string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned non-201 status: %d", resp.StatusCode)
	}

	return nil
}
