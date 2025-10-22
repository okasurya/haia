.PHONY: build run clean

# Variables
BINARY_NAME=k8sagent
GO=go

# Build the application
build:
	$(GO) build -o $(BINARY_NAME)

# Run the application in interactive mode
run: build
	./$(BINARY_NAME)

# Run the application with a specific question
run-question: build
	./$(BINARY_NAME) "$(question)"

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

# Help command to show available targets
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application in interactive mode"
	@echo "  make run-question - Run with a specific question (e.g., make run-question question=\"What pods are running?\")"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make help         - Show this help message" 