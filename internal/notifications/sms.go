// internal/notifications/sms.go
package notifications

import (
	"bytes"
	"context"
	"ecommerce-service/internal/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type smsRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}

type smsResponse struct {
	SMSMessageData struct {
		Message    string `json:"Message"`
		Recipients []struct {
			Status        int    `json:"status"`
			StatusMessage string `json:"statusMessage"`
			Number        string `json:"number"`
			MessageID     string `json:"messageId"`
			Cost          string `json:"cost"`
		} `json:"Recipients"`
	} `json:"SMSMessageData"`
}

func (s *AfricaTalkingSMSSender) Send(ctx context.Context, phone, message string) error {
	url := "https://api.africastalking.com/version1/messaging"

	payload := smsRequest{
		To:      phone,
		Message: message,
		From:    s.cfg.SenderID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", s.cfg.APIKey)
	req.Header.Set("User-Agent", "Go-AfricaTalking-Client")

	// Sandbox header only if enabled
	if s.cfg.Sandbox {
		req.Header.Set("sandbox", "true")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("AfricaTalking SMS response: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var smsResp smsResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check recipients status; usually status == 0 means success in AT API
	for _, recipient := range smsResp.SMSMessageData.Recipients {
		if recipient.Status != 0 {
			return fmt.Errorf("failed to send SMS to %s: %s", recipient.Number, recipient.StatusMessage)
		}
	}

	return nil
}
