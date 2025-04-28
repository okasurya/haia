package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

type DebugRequest struct {
	Question string `json:"question"`
}

type DebugResponse struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func handleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DebugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "OPENAI_API_KEY environment variable not set", http.StatusInternalServerError)
		return
	}

	client := openai.NewClient(apiKey)
	response, err := getCommand(client, req.Question)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get command: %v", err), http.StatusInternalServerError)
		return
	}

	command, err := extractCommand(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract command: %v", err), http.StatusInternalServerError)
		return
	}

	if !isCommandAllowed(command) {
		http.Error(w, fmt.Sprintf("Command not allowed: %s", command), http.StatusBadRequest)
		return
	}

	output, err := executeCommand(command)
	if err != nil {
		http.Error(w, fmt.Sprintf("Command execution failed: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("Output: ", output)

	resp := DebugResponse{
		Command: command,
		Output:  output,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
