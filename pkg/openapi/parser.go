package openapi

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-openapi/spec"
)

// APISpec represents the parsed OpenAPI specification
type APISpec struct {
	Spec     *spec.Swagger
	SpecJSON []byte // Raw JSON for context to the LLM
}

// ParseOpenAPIFile loads and parses an OpenAPI spec from a JSON file
func ParseOpenAPIFile(filePath string) (*APISpec, error) {
	// Read file
	specBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI file: %w", err)
	}

	// Parse JSON into spec
	swagger := &spec.Swagger{}
	if err := json.Unmarshal(specBytes, swagger); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	return &APISpec{
		Spec:     swagger,
		SpecJSON: specBytes,
	}, nil
}
