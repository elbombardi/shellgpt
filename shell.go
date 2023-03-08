package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func initShell() (*exec.Cmd, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	shell := exec.Command("/bin/sh")
	stdin, err := shell.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stdout, err := shell.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	stderr, err := shell.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := shell.Start(); err != nil {
		return nil, nil, nil, nil, err
	}
	return shell, stdin, stdout, stderr, nil
}

func runShellCommand(command string, printOnly bool) {
	// If the command has a side effect, user confirmation is needed:
	confirmationNeeded := !strings.Contains(command, "[CONFIRMATION_NOT_NEEDED]")
	command = strings.ReplaceAll(command, "[CONFIRMATION_NOT_NEEDED]", "")
	command = strings.TrimSpace(command)
	if printOnly {
		fmt.Println(command)
		return
	}
	fmt.Printf("\033[34mGenerated Command: \033[1m\033[30m%s\033[0m\n", command)
	// Prompt user to confirm whether or not to execute the command
	if confirmationNeeded && !userConfirm(command) {
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
