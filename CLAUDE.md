# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DevSpace MCP Server - A Go-based Model Context Protocol (MCP) server that wraps the DevSpace CLI, enabling AI assistants to interact with DevSpace for Kubernetes development workflows.

## Build & Development Commands

```bash
make              # Show all available targets
make build        # Build the binary
make test         # Run tests with race detection
make test-verbose # Run tests with verbose output
make fmt          # Format code
make vet          # Run go vet
make lint         # Run golangci-lint (requires installation)
make check        # Run fmt-check, vet, and test (CI pipeline)
make clean        # Remove build artifacts
```

Run a single test:
```bash
go test -v -run TestValidateCommandName ./tools/
```

Test the MCP server:
```bash
make mcp-test       # Send initialize request
make mcp-list-tools # List available tools
make mcp-inspector  # Interactive MCP inspector (requires npx)
```

## Architecture

### Core Components

- **main.go**: Entry point, creates MCP server via `mcp-go` SDK and registers all tools
- **executor/**: Command execution wrapper that shells out to `devspace` CLI
- **tools/**: MCP tool definitions and handlers

### Tool Pattern

Each tool follows a consistent pattern with two functions:

1. `Devspace<Name>Tool()` - Returns `mcp.Tool` definition with parameters
2. `Devspace<Name>Handler()` - Handles requests, calls executor, returns results

Tools are registered in `tools/tools.go` via `RegisterAll()`.

### Executor Package

The `executor` package wraps `exec.Command` for running `devspace` commands:
- `Execute()` - Default 2-minute timeout
- `ExecuteInDir()` - With working directory
- `ExecuteWithOptions()` - Custom timeout and working directory
- Returns `Result` struct with Stdout, Stderr, ExitCode, Error

### Input Validation

`tools/validate.go` provides basic validation:
- `ValidateStringParam()` - Prevents flag injection (values starting with `-`)
- `ValidateCommandName()` - Ensures command names contain only safe characters (alphanumeric, hyphen, underscore, colon)

## Key Constraints

- Interactive commands (like `devspace dev`) are not exposed - they require a terminal
- Most commands require a `devspace.yaml` file; use `working_dir` parameter to specify project location
- Default timeout: 2 minutes; build/deploy: 10 minutes
- Requires DevSpace CLI installed and in PATH
