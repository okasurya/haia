package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func getCommand(client *openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oLatest,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are an AI agent specialized in Kubernetes debugging and troubleshooting. 
Your role is to help diagnose and resolve Kubernetes cluster issues by suggesting safe kubectl commands.
Only suggest read-only commands that inspect cluster state - never suggest commands that modify resources.
Always wrap your suggested commands in <code>...</code> tags. Do not include any other text in your response.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func extractCommand(response string) (string, error) {
	start := strings.Index(response, "<code>")
	end := strings.LastIndex(response, "</code>")
	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("no command found in response")
	}
	return strings.TrimSpace(response[start+6 : end]), nil
}

func newOpenAIClient() (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	return openai.NewClient(apiKey), nil
}
