package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// APISpec represents the parsed OpenAPI specification
type APISpec struct {
	Spec     *openapi3.T
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
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(specBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Validate the spec has paths
	if doc.Paths == nil || len(doc.Paths.Map()) == 0 {
		return nil, fmt.Errorf("OpenAPI spec contains no paths")
	}

	// Pretty print the JSON for better readability in prompts
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, specBytes, "", "  ")
	if err != nil {
		// Fall back to original JSON if prettification fails
		prettyJSON.Write(specBytes)
	}

	return &APISpec{
		Spec:     doc,
		SpecJSON: prettyJSON.Bytes(),
	}, nil
}

// // APISpec represents the parsed OpenAPI specification
// type APISpec struct {
// 	Spec     *spec.Swagger
// 	SpecJSON []byte // Raw JSON for context to the LLM
// }
//
// // ParseOpenAPIFile loads and parses an OpenAPI spec from a JSON file
// func ParseOpenAPIFile(filePath string) (*APISpec, error) {
// 	// Read file
// 	specBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read OpenAPI file: %w", err)
// 	}
//
// 	// Parse JSON into spec
// 	loader := openapi3.NewLoader()
// 	doc, err := loader.LoadFromData(specBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
// 	}
//
// 	// Validate the spec has paths
// 	if doc.Paths == nil || len(doc.Paths.Map()) == 0 {
// 		return nil, fmt.Errorf("OpenAPI spec contains no paths")
// 	}
//
// 	// Pretty print the JSON for better readability in prompts
// 	var prettyJSON bytes.Buffer
// 	err = json.Indent(&prettyJSON, specBytes, "", "  ")
// 	if err != nil {
// 		// Fall back to original JSON if prettification fails
// 		prettyJSON.Write(specBytes)
// 	}
//
// 	return &APISpec{
// 		Spec:     doc,
// 		SpecJSON: prettyJSON.Bytes(),
// 	}, nil
// }
