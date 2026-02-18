package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NTFYClient sends notifications to an NTFY server
type NTFYClient struct {
	url    string
	topic  string
	client *http.Client
}

// NewNTFYClient creates a new NTFY client
func NewNTFYClient(url, topic string) *NTFYClient {
	return &NTFYClient{
		url:    url,
		topic:  topic,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Message represents an NTFY notification message
type Message struct {
	Topic    string   `json:"topic"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Priority int      `json:"priority,omitempty"` // 1=min, 3=default, 5=max
	Tags     []string `json:"tags,omitempty"`
	Click    string   `json:"click,omitempty"` // URL to open on click
}

// Send sends a notification to NTFY
func (c *NTFYClient) Send(ctx context.Context, msg *Message) error {
	// Set the topic from client config if not provided
	if msg.Topic == "" {
		msg.Topic = c.topic
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("NTFY returned status %d", resp.StatusCode)
	}

	return nil
}

// SendSimple sends a simple text notification
func (c *NTFYClient) SendSimple(ctx context.Context, title, message string) error {
	return c.Send(ctx, &Message{
		Title:   title,
		Message: message,
	})
}
