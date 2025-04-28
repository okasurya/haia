package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var allowedCommands = []string{
	"kubectl", "k", "grep", "cat", "tail", "head", "less",
}

func isCommandAllowed(cmd string) bool {
	for _, allowed := range allowedCommands {
		if strings.HasPrefix(cmd, allowed) {
			return true
		}
	}
	return false
}

func executeCommand(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	http.HandleFunc("/debug", handleDebug)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
