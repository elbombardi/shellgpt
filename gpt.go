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
// Allow user to type linux commands, correct them if needed
// If the user puts a path or a file name, do not alter it, keep it as it is.
// If no valid command can be found, give exactly the following output '(no valid command)'
// Make sure to use sudo whenever superuser previlege is needed.

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
	command, err := getGPTResponse(client, ctx, question)
	if err != nil {
		return "", err
	}
	if strings.Contains(command, "(no valid command)") {
		return "", fmt.Errorf("No valid command matches this input. Can you rephrase it ?")
	}
	return command, nil
}

func getGPTResponse(client gpt3.Client, ctx context.Context, question string) (string, error) {
	var response bytes.Buffer
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{question},
		// MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0.7),
	})
	if err != nil {
		return "", err
	}
	for _, choice := range resp.Choices {
		response.WriteString(choice.Text)
	}
	return strings.TrimSpace(response.String()), nil
}
