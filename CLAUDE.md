# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based personal LLM agent with tool calling capabilities, allowing Claude to access and modify files outside of its context window. The agent uses the Anthropic SDK for Go and provides an interactive chat interface.

## Development Commands

### Building and Running
- `make build` - Build the project
- `make install` - Install the project 
- `make test` - Run the test suite (`go test -v ./...`)
- `make run` - Run the application (executes `go run main.go`)
- `make all` - Run the complete pipeline: format, lint, install, docs, test

### Code Quality
- `make format` - Format code using `go fmt`
- `make lint` - Run golangci-lint for code linting
- `make tools` - Install required development tools (Go and golangci-lint)

### Dependency Management
- `make update` - Update dependencies with `go get -u` and `go mod tidy`

### Cleanup
- `make clean` - Remove local cache and binary directories

## Architecture

### Core Components

**Agent Structure** (`main.go:46-50`)
- `Agent` struct contains the Anthropic client, user input function, and available tools
- Manages conversation state and tool execution flow

**Tool System** (`main.go:140-145`)
- `ToolDefinition` struct defines tool name, description, JSON schema, and execution function
- Uses `jsonschema` reflection to automatically generate tool schemas from Go structs
- Three built-in tools: `read_file`, `list_files`, `edit_file`

**Conversation Loop** (`main.go:52-95`)
- Interactive chat interface with colored output (blue for user, yellow for Claude, green for tools)
- Maintains conversation history as `[]anthropic.MessageParam`
- Handles tool use responses and continues conversation until manual exit

### Built-in Tools

1. **read_file** (`main.go:163-188`)
   - Reads file contents from relative paths
   - Input: `{path: string}`

2. **list_files** (`main.go:192-253`)
   - Lists files and directories at given path (defaults to current directory)
   - Skips `.git` and `.local` directories
   - Input: `{path?: string}`

3. **edit_file** (`main.go:257-326`)
   - Performs string replacement in files
   - Creates new files if they don't exist (when old_str is empty)
   - Input: `{path: string, old_str: string, new_str: string}`

## Key Dependencies

- `github.com/anthropics/anthropic-sdk-go` - Anthropic API client
- `github.com/invopop/jsonschema` - JSON schema generation for tool definitions

## Development Environment

The Makefile handles tool installation automatically:
- Go 1.24.5 is downloaded and installed locally
- golangci-lint 2.3.0 is downloaded and installed locally
- Tools are installed to `.local/bin/` to avoid system-wide installations

## Model Configuration

Currently configured to use `anthropic.ModelClaude3_5HaikuLatest` with a 1024 token limit (`main.go:132-133`).

## Testing

The project includes comprehensive tests in `main_test.go` covering:
- **Tool Functions**: All three built-in tools (`read_file`, `list_files`, `edit_file`) with success and error cases
- **Schema Generation**: JSON schema generation for tool input validation
- **Tool Definitions**: Validation of tool metadata and structure
- **Agent Creation**: Basic agent instantiation and configuration

Run tests with `make test` to ensure all functionality works correctly.