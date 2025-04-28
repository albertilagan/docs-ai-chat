package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/albertilagan/docs-ai-chat/pkg/llm"
	"github.com/albertilagan/docs-ai-chat/pkg/openapi"
	"github.com/albertilagan/docs-ai-chat/pkg/prompt"
)

func CLIChat(apiSpec *openapi.APISpec, model string) error {
	client := llm.NewOpenRouterClient(model)
	scanner := bufio.NewScanner(os.Stdin)

	// Keep conversation history
	var conversations []llm.ChatMessage

	// Initial system message
	apiInfo := fmt.Sprintf("%s - %s", apiSpec.Spec.Info.Title, apiSpec.Spec.Info.Version)
	conversations = append(conversations, llm.ChatMessage{
		Role:    "system",
		Content: fmt.Sprintf("You are answering questions about the %s API.", apiInfo),
	})

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

		// Add user question to conversation history
		conversations = append(conversations, llm.ChatMessage{
			Role:    "user",
			Content: question,
		})

		// Build prompt with API context + conversation history
		messages := prompt.BuildOpenAPIPrompt(apiSpec, conversations)

		// Create a spinner animation
		spinnerDone := make(chan bool)
		go func() {
			spinChars := []string{"|", "/", "-", "\\"}
			i := 0
			for {
				select {
				case <-spinnerDone:
					fmt.Print("\r                    \r") // Clear the spinner
					return
				default:
					fmt.Printf("\r%s Thinking... ", spinChars[i])
					i = (i + 1) % len(spinChars)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		// Set up channels for streaming
		outputChan := make(chan string)
		doneChan := make(chan bool)

		// Start streaming in a goroutine
		go client.ChatCompletionStream(messages, outputChan, doneChan)

		// Buffer to collect the full response
		var fullResponse strings.Builder

		// Process incoming chunks
		isDone := false
		gotFirstChunk := false
		for !isDone {
			select {
			case chunk := <-outputChan:
				if !gotFirstChunk {
					// Stop the spinner when we get the first chunk
					spinnerDone <- true
					gotFirstChunk = true
				}
				fmt.Print(chunk) // Print chunk immediately
				fullResponse.WriteString(chunk)
			case <-doneChan:
				isDone = true
				if !gotFirstChunk {
					// If we never got a chunk but we're done, stop the spinner
					spinnerDone <- true
				}
			}
		}

		// Add assistant response to conversation history
		conversations = append(conversations, llm.ChatMessage{
			Role:    "assistant",
			Content: fullResponse.String(),
		})

		// Trim conversation if it gets too long (keep last 10 messages)
		if len(conversations) > 11 { // 1 system + 10 conversation turns
			// Keep the system message and latest exchanges
			conversations = append(conversations[:1], conversations[len(conversations)-10:]...)
		}

		fmt.Println() // Add a newline at the end
	}

	return nil
}
