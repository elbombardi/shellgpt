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
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize shell :
	// shell, _, _, _, err := initShell()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(2)
	// }
	// defer shell.Process.Kill()

	// Set up OpenAI client
	client, ctx, err := InitializeGPT()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	// Parse input from arguments
	input := strings.Join(os.Args[1:], " ")
	input = strings.TrimSpace(input)

	if input != "" {
		adhocMode(input, osName, ctx, client)
	} else {
		fmt.Println("Type 'exit' to terminate.")
		replMode(osName, ctx, client)
	}
}

func adhocMode(input string, osName string, ctx context.Context, client gpt3.Client) {
	// Prompt preparation :
	prompt := prepareGPTPrompt(input, osName)

	// Request ChatGPT
	command, err := getGPTResponse(client, ctx, gpt3.TextDavinci003Engine, prompt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Run the shell command
	runShellCommand(command, true)
}

func replMode(osName string, ctx context.Context, client gpt3.Client) {

	// Setup readline
	rl, err := readline.New("$ ")
	if err != nil {
		panic(err)
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

		// Add the input to the history
		history = append(history, input)

		// Prompt preparation :
		gptPrompt := prepareGPTPrompt(input, osName)

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
		runShellCommand(command, false)
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
