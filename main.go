package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
)

func getResponses(client gpt3.Client, ctx context.Context, engine string, quesiton string) (response string, err error) {
	err = client.CompletionStreamWithEngine(ctx, engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		response += resp.Choices[0].Text
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")

	return strings.TrimSpace(response), nil
}

// type NullWriter int

// func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// func chat(client gpt3.Client, ctx context.Context, humanInput string) {
// 	if humanInput != "" {
// 		conversationHistory = append(conversationHistory, fmt.Sprintf("Student: %v\n", humanInput))
// 	}

// 	prompt := situations[situation]
// 	window := min(1000, len(conversationHistory))
// 	head := len(conversationHistory) - window
// 	for _, historyentry := range conversationHistory[head:] {
// 		prompt += historyentry
// 	}
// 	robotOutput := GetResponses(client, ctx, gpt3.TextDavinci003Engine, prompt)
// 	conversationHistory = append(conversationHistory, fmt.Sprintf("%v\n###\n", robotOutput))
// 	//fmt.Println("----------------------------------------------\n", len(prompt), "\n",
// 	//		len(conversationHistory[head:]), "/", len(conversationHistory), "++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
// }

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Set up OpenAI client
	// log.SetOutput(new(NullWriter))
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("Missing API KEY. \nGo to platform.openai.com and create an API key, then store it in the environement variable OPENAI_API_KEY.")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	genericPrompt := `// Generate  valide executable linux bash shell commands for the following user input written in natural language.
// When the user input asks for exiting or terminating, output exactly the following word [BYE]. 
[user input]: {{user_input}}
[shell command]: `

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Clean up input by removing newline character and any leading/trailing whitespace
		input = strings.TrimSpace(input)

		// Prompt preparation :
		prompt := strings.ReplaceAll(genericPrompt, "{{user_input}}", input)
		fmt.Print(prompt)

		// Request ChatGPT
		output, err := getResponses(client, ctx, gpt3.TextDavinci003Engine, prompt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		// Break loop if user inputs "exit"
		if strings.Contains(input, "[BYE]") {
			fmt.Println("Bye!")
			break
		}
		fmt.Println(" output : ", output)
	}
}
