# docs-ai-chat

A command-line interface for interacting with OpenAPI specifications using OpenRouter.

## Installation

### Using Homebrew (macOS and Linux)

```bash
# Add the tap
brew tap albertilagan/tap

# Install the application
brew install daic
```

## Usage

```bash
go run cmd/main.go -spec <path-to-openapi-spec>
go run cmd/main.go -spec <path-to-openapi-spec> -model <model-name>
```

## Example

```bash
go run cmd/main.go -spec my-api.json
go run cmd/main.go -spec my-api.json -model anthropic/claude-3-sonnet
```
