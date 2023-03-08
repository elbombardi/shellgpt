package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
)

var displayCommand = strings.Join([]string{
	"ls",
	"echo",
	"pwd",
	"cd",
	"dir",
	"df",
	"du",
	"grep",
	"head",
	"lsof",
}, ", ")

const promptTemplate string = `// Generate a valide executable {{OS}} bash shell commands that matches the following natural language user input .
// if the linux commands is in this list : ({{display_commands}}), then end with '[CONFIRMATION_NOT_NEEDED]'.
[user input]: {{user_input}}
[shell command]: `

func InitializeGPT() (gpt3.Client, context.Context, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, nil, errors.New("Missing API KEY. \nGo to platform.openai.com and create an API key, then store it in the environement variable OPENAI_API_KEY.")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	return client, ctx, nil
}

func prepareGPTPrompt(userInput string, osName string) string {
	prompt := strings.Replace(promptTemplate, "{{user_input}}", userInput, 1)
	prompt = strings.Replace(prompt, "{{OS}}", osName, 1)
	prompt = strings.Replace(prompt, "{{display_commands}}", displayCommand, 1)
	fmt.Println(prompt)
	return prompt
}

func getGPTResponse(client gpt3.Client, ctx context.Context, engine string, question string) (string, error) {
	var response bytes.Buffer
	err := client.CompletionStreamWithEngine(ctx, engine, gpt3.CompletionRequest{
		Prompt:      []string{question},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0.7),
	}, func(resp *gpt3.CompletionResponse) {
		response.WriteString(resp.Choices[0].Text)
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response.String()), nil
}
