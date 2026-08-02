package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dezzles-apps/build-engine/model"
)

type EventService struct {
	host string
	port string
}

func (s *EventService) Initialise(config *model.EventServiceConfig) error {
	s.host = config.Host
	s.port = config.Port

	return nil
}

func (s *EventService) SendEvent(event model.BuildEvent) error {
	// Marshal the event to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Build the URL
	url := fmt.Sprintf("http://%s:%s/api/v1/notify", s.host, s.port)

	// Create the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Try to parse the error response
		var errorResponse map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			if errMsg, ok := errorResponse["error"]; ok {
				return fmt.Errorf("notify service error: %s", errMsg)
			}
		}
		return fmt.Errorf("notify service returned status code %d", resp.StatusCode)
	}

	return nil
}