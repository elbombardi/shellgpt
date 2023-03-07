package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/chzyer/readline"
)

// If the generated command has a side effect, end with '[CONFIRMATION_NEEDED]'.
var genericPrompt string = `// Generate a valide executable fedora linux bash shell commands that matches the following natural language user input .
[user input]: {{user_input}}
[shell command]: `

func main() {
	// Parse input from arguments
	input := strings.Join(os.Args[1:], " ")
	input = strings.TrimSpace(input)
	//fmt.Printf("'%s'\n", input)

	// Set up OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing API KEY. \nGo to platform.openai.com and create an API key, then store it in the environement variable OPENAI_API_KEY.")
		os.Exit(0)

	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	if input != "" {
		adhocMode(input, ctx, client)
	} else {
		replMode(ctx, client)
	}
}

func adhocMode(input string, ctx context.Context, client gpt3.Client) {
	// Prompt preparation :
	// Prompt preparation :
	prompt := strings.ReplaceAll(genericPrompt, "{{user_input}}", input)

	// Request ChatGPT
	command, err := getResponses(client, ctx, gpt3.TextDavinci003Engine, prompt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(command)
}

func getResponses(client gpt3.Client, ctx context.Context, engine string, question string) (string, error) {
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

func replMode(ctx context.Context, client gpt3.Client) {
	printIntro()

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
		prompt := strings.ReplaceAll(genericPrompt, "{{user_input}}", input)

		// Request ChatGPT
		command, err := getResponses(client, ctx, gpt3.TextDavinci003Engine, prompt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		// Break loop if user inputs "exit"
		if input == "exit" {
			fmt.Println("Bye!")
			break
		}

		// If the command has a side effect, user confirmation is needed:
		confirmationNeeded := strings.Contains(command, "[CONFIRMATION_NEEDED]")
		command = strings.ReplaceAll(command, "[CONFIRMATION_NEEDED]", "")
		command = strings.TrimSpace(command)
		fmt.Printf("Generated shell command: %s\n", command)

		// Prompt user to confirm whether or not to execute the command
		if confirmationNeeded && !userConfirm(command) {
			continue
		}

		// Run the shell command
		runShellCommand(command)
	}
}

func printIntro() {
	fmt.Println("Type 'exit' to terminate.")
}

func userConfirm(command string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Execute command? (y/N): ")
	executionConfirmation, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	executionConfirmation = strings.TrimSpace(strings.ToLower(executionConfirmation))
	return executionConfirmation == "y" || executionConfirmation == "yes"

}

func runShellCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
	}
}
