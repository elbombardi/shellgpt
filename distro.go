package main

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func osName() (string, error) {
	// Check if OS is linux
	if runtime.GOOS != "linux" {
		return "", errors.New("Only Linux OS is supported.")
	}

	// Check if "/etc/os-release" exists
	if _, err := os.Stat("/etc/os-release"); err == nil {
		// Read the contents of "/etc/os-release"
		data, err := ioutil.ReadFile("/etc/os-release")
		if err != nil {
			return "", err
		}

		// Find the value of the "PRETTY_NAME" field
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				return strings.Trim(line[len("PRETTY_NAME="):], `"`), nil
			}
		}
	}

	// Check if "/etc/system-release" exists
	if _, err := os.Stat("/etc/system-release"); err == nil {
		// Read the contents of "/etc/system-release"
		data, err := ioutil.ReadFile("/etc/system-release")
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(data)), nil
	}

	return "linux", nil
}
