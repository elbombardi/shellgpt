package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/chzyer/readline"
)

func main() {
	// Detecting linux distribution :
	osName, err := osName()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Set up OpenAI client
	client, ctx, err := InitializeGPT()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// Parse input from arguments
	fmt.Println("Type 'exit' to terminate.")
	replMode(osName, ctx, client)
}

func replMode(osName string, ctx context.Context, client gpt3.Client) {

	// Setup readline
	rl, err := readline.New("$ ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rl.Close()

	history := []string{}
	for {
		// Read user input
		input, err := rl.Readline()
		if err != nil {
			fmt.Println("Bye!")
			break
		}

		// Clean up input by removing newline character and any leading/trailing whitespace
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Prompt preparation :
		gptPrompt := prepareGPTPrompt(input, osName, history)

		// Add the input to the history
		history = append(history, input)

		// Request ChatGPT
		command, err := getGPTResponse(client, ctx, gpt3.TextDavinci003Engine, gptPrompt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		// Break loop if user inputs "exit"
		if input == "exit" {
			fmt.Println("\n\033[34mBye!\033[0m")
			break
		}

		// Run the shell command
		runShellCommand(command)
	}
}

func userConfirm(command string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\033[34mExecute command? (y/N): \033[0m")
	executionConfirmation, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	executionConfirmation = strings.TrimSpace(strings.ToLower(executionConfirmation))
	return executionConfirmation == "y" || executionConfirmation == "yes"
}
