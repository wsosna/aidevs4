# aidevs4

My solutions to [AI Devs 4](https://www.aidevs.pl/) exercises, written in Go. Exercise content is confidential, but
below I describe the concepts each solution explores.

## Notable Implementations

**Prompt caching** — `internal/ai/client.go` hashes each prompt (SHA-256) and stores the response in a local file.
Subsequent identical prompts are served from cache, skipping the API call entirely.

## Exercises

### W1E1 — Structured Output

Classifies people from a CSV using
OpenAI's [Structured Output](https://platform.openai.com/docs/guides/structured-outputs)
feature.                                                                                                                                                               
See `internal/exercises/e0101.go`