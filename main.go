package main

import (
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

	history := []*HistoryEntry{}
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

		// Break loop if user inputs "exit"
		if input == "exit" {
			fmt.Println("\n\033[34mBye!\033[0m")
			break
		}

		// Clear the screen and remove history if the user inputs "clear"
		if input == "clear" || input == "cls" {
			fmt.Print("\033[H\033[2J")
			continue
		}

		// Prompt preparation :
		gptPrompt := prepareGPTPrompt(input, osName, history)

		// Request ChatGPT
		command, err := generateCommand(client, ctx, gptPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%v\033[0m\n", err)
			continue
		}

		// Run the shell command
		cancelled := runShellCommand(command, history)

		// Add the input to the history
		if !cancelled {
			history = append(history, &HistoryEntry{
				userInput: input,
				command:   command,
			})
		}

	}
}
