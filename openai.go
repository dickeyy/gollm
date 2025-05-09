package gollm

import (
	"context"
	"errors"
	"fmt"

	oai "github.com/sashabaranov/go-openai"
)

// OpenAIModel implements the gollm.LLM interface for OpenAI models.
type OpenAIModel struct {
	modelID string
	apiKey  string
	apiURL  string
	client  *oai.Client
}

var defaultAPIURL = "https://api.openai.com/v1"

// NewOpenAiModel creates a new instance of OpenAIModel.
// This is the function that will be registered with the factory.
func NewOpenAiModel(modelName string, apiKey string, apiURL string) (*OpenAIModel, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	config := oai.DefaultConfig(apiKey)
	if apiURL != "" {
		config.BaseURL = apiURL
	} else {
		config.BaseURL = defaultAPIURL
	}

	client := oai.NewClientWithConfig(config)

	return &OpenAIModel{
		modelID: modelName,
		apiKey:  apiKey,
		apiURL:  config.BaseURL,
		client:  client,
	}, nil
}

// Chat implements the gollm.LLM interface.
func (m *OpenAIModel) Chat(structure ChatStructure) (*ChatResponse, error) {
	if m.client == nil {
		return nil, errors.New("OpenAI client not initialized")
	}

	var oaiMessages []oai.ChatCompletionMessage
	for _, msg := range structure.Messages {
		oaiMessages = append(oaiMessages, oai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	if len(oaiMessages) == 0 {
		return nil, errors.New("no messages provided in ChatStructure")
	}

	resp, err := m.client.CreateChatCompletion(
		context.Background(),
		oai.ChatCompletionRequest{
			Model:    m.modelID,
			Messages: oaiMessages,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return nil, errors.New("no response content received from OpenAI")
	}

	return &ChatResponse{
		Text: resp.Choices[0].Message.Content,
	}, nil
}

// init registers the OpenAI model constructors with the gollm factory.
func init() {
	constructor := func(modelName string, apiKey string, apiURL string) (LLM, error) {
		model, err := NewOpenAiModel(modelName, apiKey, apiURL)
		if err != nil {
			return nil, err
		}
		return model, nil
	}

	modelsToRegister := []string{"gpt-4o", "gpt-4", "gpt-3.5-turbo", "gpt-4o-mini"}
	for _, modelName := range modelsToRegister {
		RegisterModel(modelName, constructor)
	}
}
