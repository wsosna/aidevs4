package ai

import (
	"aidevs4/internal/cache"
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const defaultModel = openai.ChatModelGPT4oMini

type Client struct {
	openai openai.Client
	cache  *cache.LocalCache
}

func NewClient() *Client {
	return &Client{
		openai: openai.NewClient(),
		cache:  cache.DefaultLocalCache,
	}
}

type requestOptions struct {
	model     openai.ChatModel
	format    responses.ResponseFormatTextConfigUnionParam
	hasFormat bool
}

type RequestOption func(*requestOptions)

func WithModel(model openai.ChatModel) RequestOption {
	return func(o *requestOptions) {
		o.model = model
	}
}

func WithFormat(schema map[string]any) RequestOption {
	return func(o *requestOptions) {
		o.format = responses.ResponseFormatTextConfigParamOfJSONSchema("classification_result", schema)
		o.hasFormat = true
	}
}

func (c *Client) Request(ctx context.Context, prompt string, opts ...RequestOption) (string, error) {
	if cached, ok := c.cache.Get(prompt); ok {
		fmt.Printf("Request cached - skip OpenAI API call")
		return cached, nil
	}

	o := &requestOptions{
		model: defaultModel,
	}
	for _, opt := range opts {
		opt(o)
	}

	params := responses.ResponseNewParams{
		Model: o.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
	}
	if o.hasFormat {
		params.Text = responses.ResponseTextConfigParam{Format: o.format}
	}

	resp, err := c.openai.Responses.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}

	result := resp.OutputText()
	if err := c.cache.Set(prompt, result); err != nil {
		return "", fmt.Errorf("failed to cache response: %w", err)
	}

	return result, nil
}
