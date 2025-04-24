package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/albertilagan/docs-ai-chat/pkg/llm"
	"github.com/albertilagan/docs-ai-chat/pkg/openapi"
	"github.com/albertilagan/docs-ai-chat/pkg/prompt"
)

// CLIChat runs a command-line chat interface
func CLIChat(apiSpec *openapi.APISpec, model string) error {
	client := llm.NewOpenRouterClient(model)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("OpenAPI Assistant ready! Ask questions about the API (type 'exit' to quit):")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		question := scanner.Text()
		if strings.ToLower(question) == "exit" {
			break
		}

		messages := prompt.BuildOpenAPIPrompt(apiSpec, question)

		fmt.Println("\nThinking...")
		response, err := client.ChatCompletion(messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("\n" + response)
	}

	return nil
}
