package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runShellCommand(command string) {
	// If the command has a side effect, user confirmation is needed:
	command = strings.ReplaceAll(command, "[CONFIRMATION_NOT_NEEDED]", "")
	command = strings.TrimSpace(command)
	fmt.Printf("\033[34mGenerated Command: \033[7m%s\033[0m\n", command)

	// Prompt user to confirm whether or not to execute the command
	if !userConfirm(command) {
		fmt.Println("\n\033[34m---------------\nCancelled!\n\033[0m")
		return
	}
	fmt.Printf("\033[34m---------------\033[0m\n")

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	fmt.Println("\033[34m---------------")
	if err != nil {
		fmt.Printf("Failed to execute command: %v\033[0m\n", err)
		return
	}

	fmt.Println("\033[34mDone!\033[0m")
}
