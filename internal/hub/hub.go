package hub

import (
	"fmt"
	"os"
	"strings"
)

const urlPlaceholder = "tutaj-twój-klucz"

// VerifyAnswer sends task and answer to the hub server for verification.
// Requires HUB_API_KEY and HUB_SERVER_URL environment variables.
// Returns the response message body on success.
func VerifyAnswer(task string, answer any) (string, error) {
	apiKey := os.Getenv("HUB_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HUB_API_KEY environment variable is not set")
	}
	serverURL := os.Getenv("HUB_SERVER_URL")
	if serverURL == "" {
		return "", fmt.Errorf("HUB_SERVER_URL environment variable is not set")
	}

	return sendRequest(serverURL+"/verify", apiKey, task, answer)
}

// FetchFile fetches content from url (replacing urlPlaceholder with HUB_API_KEY)
// and returns it as a string.
func FetchFile(url string) (string, error) {
	apiKey := os.Getenv("HUB_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HUB_API_KEY environment variable is not set")
	}

	resolvedURL := strings.ReplaceAll(url, urlPlaceholder, apiKey)
	return fetchContent(resolvedURL)
}
