package prompt

import (
	"fmt"
	"github.com/albertilagan/docs-ai-chat/pkg/llm"
	"github.com/albertilagan/docs-ai-chat/pkg/openapi"
)

func BuildOpenAPIPrompt(apiSpec *openapi.APISpec, conversations []llm.ChatMessage) []llm.ChatMessage {
	// The first message should always be the system prompt
	systemPrompt := conversations[0].Content

	// Add API context to the system prompt
	systemPromptWithContext := fmt.Sprintf(`%s

Here is the OpenAPI specification for this API:

%s

Answer questions based on this specification. Be specific about endpoints, parameters, request formats, and response structures.`,
		systemPrompt,
		string(apiSpec.SpecJSON))

	// Create a new array with the enhanced system prompt and user's conversation
	result := make([]llm.ChatMessage, len(conversations))
	result[0] = llm.ChatMessage{Role: "system", Content: systemPromptWithContext}

	// Copy the rest of the conversation
	for i := 1; i < len(conversations); i++ {
		result[i] = conversations[i]
	}

	return result
}
