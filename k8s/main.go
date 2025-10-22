package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var allowedCommands = []string{
	"kubectl", "k", "grep", "cat", "tail", "head", "less",
}

func isCommandAllowed(cmd string) bool {
	// Split the command by pipe and check each part
	parts := strings.Split(cmd, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		cmdParts := strings.Fields(part)
		if len(cmdParts) == 0 {
			continue
		}
		allowed := false
		for _, allowedCmd := range allowedCommands {
			if cmdParts[0] == allowedCmd {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

func executeCommand(command string) (string, error) {
	// Use shell to execute the command to handle pipes
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func confirmExecution(command string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to execute: %s? (y/N): ", command)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading confirmation: %v", err)
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func processQuestion(question string) {
	client, err := newOpenAIClient()
	if err != nil {
		log.Printf("Failed to create OpenAI client: %v", err)
		return
	}

	response, err := getCommand(client, question)
	if err != nil {
		log.Printf("Failed to get command: %v", err)
		return
	}

	command, err := extractCommand(response)
	if err != nil {
		log.Printf("Failed to extract command: %v", err)
		return
	}

	if !isCommandAllowed(command) {
		log.Printf("Command not allowed: %s", command)
		return
	}

	if !confirmExecution(command) {
		fmt.Println("Command execution cancelled.")
		return
	}

	log.Printf("Executing command: %s", command)
	output, err := executeCommand(command)
	if err != nil {
		log.Printf("Command execution failed: %v", err)
		return
	}

	fmt.Println(output)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [question]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s \"What pods are running in the default namespace?\"\n", os.Args[0])
	}
	flag.Parse()

	// If arguments are provided, process them and exit
	if len(flag.Args()) > 0 {
		question := strings.Join(flag.Args(), " ")
		processQuestion(question)
		return
	}

	// Interactive mode
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Kubernetes Debug Assistant (type 'exit' or 'quit' to end)")
	fmt.Println("------------------------------------------------------")

	for {
		fmt.Print("> ")
		question, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		question = strings.TrimSpace(question)
		if question == "" {
			continue
		}

		if strings.ToLower(question) == "exit" || strings.ToLower(question) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		processQuestion(question)
	}
}
