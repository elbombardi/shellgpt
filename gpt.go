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

const PROMPT_TEMPLATE string = `// Generate a valide executable {{OS}} bash shell commands that matches the following natural language user input .
// If no valid command can be found, give exactly the following output '(no valid command)'

{{history_log}}
[user input]: {{user_input}}
[shell command]: `

type HistoryEntry struct {
	userInput string
	command   string
}

func (entry *HistoryEntry) String() string {
	return fmt.Sprintf("\n[user input]: %s\n[shell command]: %s", entry.userInput, entry.command)
}

func InitializeGPT() (gpt3.Client, context.Context, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, nil, errors.New("Missing API KEY. \nGo to platform.openai.com and create an API key, then store it in the environement variable OPENAI_API_KEY.")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	return client, ctx, nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func prepareGPTPrompt(userInput string, osName string, history []*HistoryEntry) string {
	prompt := strings.Replace(PROMPT_TEMPLATE, "{{user_input}}", userInput, 1)
	prompt = strings.Replace(prompt, "{{OS}}", osName, 1)
	historyLog := ""
	for _, entry := range history[max(len(history)-10, 0):] {
		historyLog += entry.String()
	}
	prompt = strings.Replace(prompt, "{{history_log}}", historyLog, 1)
	return prompt
}

func generateCommand(client gpt3.Client, ctx context.Context, question string) (string, error) {
	command, err := getGPTResponse(client, ctx, gpt3.TextDavinci003Engine, question)
	if err != nil {
		return "", err
	}
	if strings.Contains(command, "(no valid command)") {
		return "", fmt.Errorf("No valid command matches this input. Can you rephrase it ?")
	}
	return command, nil
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
