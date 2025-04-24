package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/albertilagan/docs-ai-chat/pkg/chat"
	"github.com/albertilagan/docs-ai-chat/pkg/openapi"
)

// cmd/main.go
func main() {
	// Define command-line flags
	specFile := flag.String("spec", "", "Path to OpenAPI JSON specification file")
	model := flag.String("model", "anthropic/claude-3-sonnet", "OpenRouter model to use")
	flag.Parse()

	if *specFile == "" {
		fmt.Println("Error: You must provide an OpenAPI spec file with -spec")
		os.Exit(1)
	}

	// Parse OpenAPI spec
	apiSpec, err := openapi.ParseOpenAPIFile(*specFile)
	if err != nil {
		fmt.Printf("Error parsing OpenAPI spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully loaded OpenAPI spec: %s %s\n\n",
		apiSpec.Spec.Info.Title,
		apiSpec.Spec.Info.Version)

	fmt.Printf("Using model: %s\n\n", *model)

	// Start CLI chat with the specified model
	if err := chat.CLIChat(apiSpec, *model); err != nil {
		fmt.Printf("Chat error: %v\n", err)
		os.Exit(1)
	}
}
