package hub

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildRequestBody(t *testing.T) {
	answer := []map[string]any{
		{
			"name":    "Jan",
			"surname": "Kowalski",
			"gender":  "M",
			"born":    1987,
			"city":    "Warszawa",
			"tags":    []string{"tag1", "tag2"},
		},
		{
			"name":    "Anna",
			"surname": "Nowak",
			"gender":  "F",
			"born":    1993,
			"city":    "Grudziądz",
			"tags":    []string{"tagA", "tagB", "tagC"},
		},
	}

	got, err := buildRequestBody("tutaj-twój-klucz-api", "people", answer)
	if err != nil {
		t.Fatalf("buildRequestBody returned error: %v", err)
	}

	want := `{"apikey":"tutaj-twój-klucz-api","task":"people","answer":[{"born":1987,"city":"Warszawa","gender":"M","name":"Jan","surname":"Kowalski","tags":["tag1","tag2"]},{"born":1993,"city":"Grudziądz","gender":"F","name":"Anna","surname":"Nowak","tags":["tagA","tagB","tagC"]}]}`

	if string(got) != want {
		t.Errorf("body mismatch:\ngot:  %s\nwant: %s", got, want)
	}
}

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	answer := []map[string]any{{"key": "value"}}
	msg, err := sendRequest(server.URL, "test-key", "test-task", answer)
	if err != nil {
		t.Fatalf("sendRequest returned error: %v", err)
	}
	if msg != "ok" {
		t.Errorf("expected message %q, got %q", "ok", msg)
	}
}
