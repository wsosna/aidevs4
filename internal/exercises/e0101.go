package exercises

import (
	"aidevs4/internal/ai"
	"aidevs4/internal/hub"
	"aidevs4/internal/xio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type PersonClassification struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Gender  string   `json:"gender"`
	Born    int      `json:"born"`
	City    string   `json:"city"`
	Tags    []string `json:"tags"`
}

type ClassificationResult struct {
	People []PersonClassification `json:"people"`
}

var classificationSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"people": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name":    map[string]any{"type": "string"},
					"surname": map[string]any{"type": "string"},
					"gender":  map[string]any{"type": "string"},
					"born":    map[string]any{"type": "integer"},
					"city":    map[string]any{"type": "string"},
					"tags": map[string]any{
						"type":  "array",
						"items": map[string]any{"type": "string"},
					},
				},
				"required":             []string{"name", "surname", "gender", "born", "city", "tags"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"people"},
	"additionalProperties": false,
}

func RunExercise1() {
	content, err := hub.FetchFile("https://hub.ag3nts.org/data/tutaj-twój-klucz/people.csv")
	if err != nil {
		panic(err)
	}
	csv, err := xio.FilterCsvFromString(content, 1986, 2006, "", "", "M", "", "Grudziądz")
	if err != nil {
		panic(err)
	}
	if csv == nil {
		panic("No records!")
	}

	fmt.Printf("Result of filtering %v rows\n", len(csv)-1)

	lines := make([]string, len(csv))
	for i, row := range csv {
		lines[i] = strings.Join(row, ",")
	}
	csvStr := strings.Join(lines, "\n")

	systemMsg := "<system>" +
		"You are given list of people in csv format where columns are splitted with coma (,). Each row contains some basic" +
		"information about the person and MOST IMPORTANTLY: their job description." +
		"Classify each of them, based on their job description with available tags:" +
		"- IT\n" +
		"- transport\n" +
		"- edukacja\n" +
		"- medycyna\n" +
		"- praca z ludźmi\n" +
		"- praca z pojazdami\n" +
		"- praca fizyczna\n" +
		"One person can have more than just one tag. Return the result as the list of JSON objects in given format." +
		"</system>" +
		"<csv>\n" +
		csvStr +
		"\n</csv>"

	ctx := context.Background()
	client := ai.NewClient()

	fmt.Printf("Calling OpenAI API\n")
	output, err := client.Request(ctx, systemMsg, ai.WithFormat(classificationSchema))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n----- AI RESPONSE: ------\n\n%s", output)

	var classification ClassificationResult
	if err := json.Unmarshal([]byte(output), &classification); err != nil {
		panic(fmt.Errorf("failed to parse AI response: %w", err))
	}

	var transport []PersonClassification
	for _, p := range classification.People {
		for _, tag := range p.Tags {
			if tag == "transport" {
				transport = append(transport, p)
				break
			}
		}
	}

	fmt.Printf("Result of filtering by transport tag %v", len(transport))

	result, err := hub.VerifyAnswer("people", transport)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n----- HUB RESPONSE: ------\n\n%s", result)
}
