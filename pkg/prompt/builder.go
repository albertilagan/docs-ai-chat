package prompt

import (
	"encoding/json"
	"fmt"
	"github.com/albertilagan/docs-ai-chat/pkg/llm"
	"github.com/albertilagan/docs-ai-chat/pkg/openapi"
)

// BuildOpenAPIPrompt creates a system prompt for the LLM with OpenAPI context
func BuildOpenAPIPrompt(apiSpec *openapi.APISpec, userQuestion string) []llm.ChatMessage {
	// Create a compact version of the spec for context
	specSummary, _ := json.Marshal(map[string]interface{}{
		"info":        apiSpec.Spec.Info,
		"paths":       apiSpec.Spec.Paths,
		"definitions": apiSpec.Spec.Definitions,
		"parameters":  apiSpec.Spec.Parameters,
		"responses":   apiSpec.Spec.Responses,
	})

	systemPrompt := fmt.Sprintf(`You are an API documentation assistant. Answer questions about this OpenAPI specification:

%s

Be concise but thorough. If asked about endpoints, include the path, method, parameters, 
request body structure, and response format. If the information isn't in the spec, say so.`,
		string(specSummary))

	return []llm.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userQuestion},
	}
}
