package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type request struct {
	APIKey string `json:"apikey"`
	Task   string `json:"task"`
	Answer any    `json:"answer"`
}


func buildRequestBody(apiKey, task string, answer any) ([]byte, error) {
	payload := request{
		APIKey: apiKey,
		Task:   task,
		Answer: answer,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	return body, nil
}

func sendRequest(serverURL, apiKey, task string, answer any) (string, error) {
	body, err := buildRequestBody(apiKey, task, answer)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(respBody), nil
}

func fetchContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
