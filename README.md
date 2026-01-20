# DevSpace MCP Server

A Model Context Protocol (MCP) server that exposes [DevSpace](https://devspace.sh/) CLI functionality as MCP tools, enabling AI assistants like Claude to interact with DevSpace for Kubernetes development workflows.

## Overview

This MCP server wraps the DevSpace CLI and provides a set of tools that allow AI assistants to:

- Query Kubernetes cluster information (namespaces, contexts, deployments)
- Build and deploy applications to Kubernetes
- Analyze namespaces for potential issues
- Retrieve pod logs
- Execute predefined DevSpace commands
- Manage DevSpace configurations and profiles

## Requirements

- **Go 1.21+** - For building from source
- **DevSpace CLI** - Must be installed and available in PATH ([Installation Guide](https://devspace.sh/docs/getting-started/installation))
- **Kubernetes cluster** - With valid kubeconfig for cluster access

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/devspace-mcp.git
cd devspace-mcp

# Build the binary
go build -o devspace-mcp .

# Optionally, move to a directory in your PATH
mv devspace-mcp /usr/local/bin/
```

### Verify Installation

```bash
# Check that devspace CLI is available
devspace version

# Test the MCP server
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./devspace-mcp
```

## Usage

### With Claude Code

Add the server to your Claude Code MCP configuration at `~/.claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "devspace": {
      "command": "/path/to/devspace-mcp"
    }
  }
}
```

Restart Claude Code to load the new MCP server. You can then ask Claude to perform DevSpace operations like:

- "List all Kubernetes contexts"
- "Deploy the application to the staging namespace"
- "Show me the logs from the web pod"
- "Analyze the current namespace for issues"

### With MCP Inspector

Test the server interactively using the MCP Inspector:

```bash
npx @anthropic-ai/mcp-inspector ./devspace-mcp
```

### Programmatic Usage

The server communicates via stdio using the JSON-RPC 2.0 protocol:

```bash
# Initialize and list tools
cat << 'EOF' | ./devspace-mcp
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
EOF
```

## Tools Reference

### devspace_version

Get the DevSpace CLI version.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| *(none)* | - | - | - |

**Example:**
```json
{"name": "devspace_version", "arguments": {}}
```

---

### devspace_list_namespaces

List Kubernetes namespaces accessible from the current context.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `kube_context` | string | No | Kubernetes context to use |

**Example:**
```json
{"name": "devspace_list_namespaces", "arguments": {"kube_context": "my-cluster"}}
```

---

### devspace_list_contexts

List all available Kubernetes contexts from kubeconfig.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| *(none)* | - | - | - |

**Example:**
```json
{"name": "devspace_list_contexts", "arguments": {}}
```

---

### devspace_list_deployments

List deployments and their current status.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace to query |
| `kube_context` | string | No | Kubernetes context to use |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_list_deployments", "arguments": {"namespace": "default"}}
```

---

### devspace_list_profiles

List available DevSpace profiles defined in `devspace.yaml`.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_list_profiles", "arguments": {}}
```

---

### devspace_list_vars

List variables defined in the active DevSpace configuration.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `profile` | string | No | Profile to use when resolving variables |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_list_vars", "arguments": {"profile": "production"}}
```

---

### devspace_print

Print the fully resolved DevSpace configuration as YAML.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `profile` | string | No | Profile to apply when resolving the configuration |
| `skip_info` | boolean | No | Only print the configuration without additional info |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_print", "arguments": {"profile": "dev", "skip_info": true}}
```

---

### devspace_analyze

Analyze a Kubernetes namespace for potential problems and issues.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace to analyze |
| `kube_context` | string | No | Kubernetes context to use |
| `wait` | boolean | No | Wait for pods to be ready before analyzing (default: true) |
| `timeout` | number | No | Timeout in seconds (default: 120, max: 600) |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_analyze", "arguments": {"namespace": "my-app", "timeout": 60}}
```

---

### devspace_logs

Get logs from a pod in the Kubernetes cluster.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace |
| `pod` | string | No | Specific pod name to get logs from |
| `container` | string | No | Container name within the pod |
| `label_selector` | string | No | Label selector to filter pods (e.g., `app=myapp`) |
| `lines` | number | No | Maximum number of lines to return (default: 200, max: 10000) |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_logs", "arguments": {"label_selector": "app=web", "lines": 100}}
```

---

### devspace_build

Build all images defined in `devspace.yaml`.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace |
| `kube_context` | string | No | Kubernetes context to use |
| `profile` | string | No | Profile to use |
| `skip_push` | boolean | No | Skip pushing images to registry |
| `tag` | string | No | Tag to use for built images |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_build", "arguments": {"profile": "dev", "skip_push": true, "tag": "v1.2.3"}}
```

---

### devspace_deploy

Deploy the project to Kubernetes using DevSpace.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace to deploy to |
| `kube_context` | string | No | Kubernetes context to use |
| `profile` | string | No | Profile to use |
| `force_build` | boolean | No | Force rebuilding images even if not changed |
| `force_deploy` | boolean | No | Force redeployment even if not changed |
| `skip_build` | boolean | No | Skip building images |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_deploy", "arguments": {"namespace": "staging", "profile": "staging", "skip_build": true}}
```

---

### devspace_purge

**WARNING: Destructive operation.** Delete all deployed Kubernetes resources for the project. This cannot be undone.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `namespace` | string | No | Kubernetes namespace |
| `kube_context` | string | No | Kubernetes context to use |
| `profile` | string | No | Profile to use |
| `force_purge` | boolean | No | Force purge even if resources are in use |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_purge", "arguments": {"namespace": "dev", "force_purge": true}}
```

---

### devspace_run

Execute a predefined command from `devspace.yaml`.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `command` | string | **Yes** | Name of the command to run (as defined in devspace.yaml) |
| `args` | string | No | Arguments to pass to the command (space-separated) |
| `working_dir` | string | No | Working directory containing devspace.yaml |

**Example:**
```json
{"name": "devspace_run", "arguments": {"command": "migrate", "args": "--force"}}
```

## Project Structure

```
devspace-mcp/
├── main.go              # Entry point, MCP server setup
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── executor/
│   ├── executor.go      # Command execution wrapper with timeout support
│   └── executor_test.go # Unit tests for executor
└── tools/
    ├── tools.go         # Tool registration
    ├── validate.go      # Input validation helpers
    ├── validate_test.go # Unit tests for validation
    ├── version.go       # devspace_version tool
    ├── list.go          # List tools (namespaces, contexts, deployments, profiles, vars)
    ├── print.go         # devspace_print tool
    ├── analyze.go       # devspace_analyze tool
    ├── logs.go          # devspace_logs tool
    ├── build.go         # devspace_build tool
    ├── deploy.go        # devspace_deploy tool
    ├── purge.go         # devspace_purge tool
    └── run.go           # devspace_run tool
```

## Limitations

- **Interactive commands not supported**: Commands like `devspace dev` that require an interactive terminal are not exposed as tools. Use `devspace deploy` + `devspace logs` workflow instead.
- **Requires devspace.yaml**: Most commands require a `devspace.yaml` configuration file. Use the `working_dir` parameter to specify the project directory if not running from the project root.
- **Cluster access required**: Commands that interact with Kubernetes require valid kubeconfig and cluster access.

## Timeouts

- Default command timeout: **2 minutes**
- Build/Deploy commands: **10 minutes**
- Analyze command: Configurable via `timeout` parameter (default: 120 seconds, max: 600 seconds)

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Projects

- [DevSpace](https://devspace.sh/) - The DevSpace CLI
- [MCP Specification](https://modelcontextprotocol.io/) - Model Context Protocol documentation
- [mcp-go](https://github.com/mark3labs/mcp-go) - Go SDK for MCP servers
