# docs-ai-chat

A command-line interface for interacting with OpenAPI specifications using OpenRouter.

## Usage

```bash
go run cmd/main.go -spec <path-to-openapi-spec> -model <model-name>
```

## Example

```bash
go run cmd/main.go -spec my-api.json -model anthropic/claude-3-sonnet
```
