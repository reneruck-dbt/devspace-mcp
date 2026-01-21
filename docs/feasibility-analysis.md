# DevSpace MCP Server - Feature Feasibility Analysis

This document analyzes the feasibility of implementing features requested in user feedback, based on DevSpace CLI capabilities.

## Summary

| # | Feature Request | Feasible | Effort | Notes |
|---|-----------------|----------|--------|-------|
| 1 | Consistent working_dir | ✅ Yes | Low | Code improvement |
| 2 | Better analyze output | ⚠️ Partial | Low | Add flags, may still be minimal when healthy |
| 3 | devspace_list_pods | ✅ Yes* | Medium | Requires kubectl wrapper |
| 4 | devspace_exec | ✅ Yes | Medium | Via `devspace enter --tty=false` |
| 5 | devspace_dev | ❌ No | N/A | Interactive, requires terminal |
| 6 | devspace_port_forward | ⚠️ Partial | Medium | List only; actual forwarding needs kubectl |
| 7 | Better error messages | ✅ Yes | Medium | Add error pattern detection |
| 8 | Streaming progress | ❌ No | N/A | MCP protocol limitation |
| 9 | devspace_status | ✅ Yes | Medium | Composite of multiple commands |
| 10 | Log filtering | ⚠️ Partial | Low-Medium | Some features, not all |

---

## Detailed Analysis

### 1. Inconsistent working_dir Requirement

**Status:** ✅ Feasible | ✅ IMPLEMENTED

**Problem:** Some tools work without `working_dir`, others fail with "Cannot find a devspace.yaml"

**Root Cause:** Commands like `version`, `list contexts`, `list namespaces` don't need a devspace.yaml, while `deploy`, `build`, `list deployments` do.

**Implementation:**

```go
// In each handler, add consistent detection
func requiresDevspaceYaml(workingDir string) error {
    configPath := filepath.Join(workingDir, "devspace.yaml")
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return fmt.Errorf("devspace.yaml not found in %s. Use working_dir parameter to specify the project location", workingDir)
    }
    return nil
}
```

**Changes:**
- Add validation function in `tools/validate.go`
- Document in tool descriptions which require `working_dir`
- Standardize error messages

---

### 2. devspace_analyze Returns Minimal Output

**Status:** ⚠️ Partially Feasible

**Problem:** Only shows "Checking status..." with no actual analysis

**Root Cause:** When everything is healthy, DevSpace's analyze output is minimal by design. It reports problems, not health.

**CLI Capabilities:**
```
--ignore-pod-restarts   Ignore restart events of running pods
--patient               Wait for all resources to be ready before reporting
--timeout int           Timeout for waiting (default 120)
--wait                  Wait for pods to get ready (default true)
```

**Implementation:**

```go
func DevspaceAnalyzeTool() mcp.Tool {
    return mcp.NewTool("devspace_analyze",
        mcp.WithDescription("Analyzes namespace for potential problems..."),
        mcp.WithString("namespace", mcp.Description("Target namespace")),
        mcp.WithString("working_dir", mcp.Description("DevSpace project directory")),
        mcp.WithBoolean("patient", mcp.Description("Wait for all resources before analyzing")),
        mcp.WithBoolean("ignore_pod_restarts", mcp.Description("Ignore pod restart events")),
        mcp.WithNumber("timeout", mcp.Description("Analysis timeout in seconds (default 120)")),
    )
}

func DevspaceAnalyzeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := []string{"analyze", "--no-colors"}

    if patient, _ := request.Params.Arguments["patient"].(bool); patient {
        args = append(args, "--patient")
    }
    if ignoreRestarts, _ := request.Params.Arguments["ignore_pod_restarts"].(bool); ignoreRestarts {
        args = append(args, "--ignore-pod-restarts")
    }
    if timeout, ok := request.Params.Arguments["timeout"].(float64); ok {
        args = append(args, fmt.Sprintf("--timeout=%d", int(timeout)))
    }
    // ... execute
}
```

**Limitation:** Output will still be minimal when cluster is healthy. Consider documenting this behavior.

---

### 3. No Pod Listing Capability

**Status:** ✅ Feasible (with kubectl) | ✅ IMPLEMENTED

**Problem:** No `devspace list pods` command exists

**Analysis:** DevSpace intentionally doesn't duplicate kubectl functionality. Pod selection is handled via label selectors in other commands.

**Implementation Options:**

**Option A: kubectl wrapper (Recommended)**
```go
func DevspaceListPodsTool() mcp.Tool {
    return mcp.NewTool("devspace_list_pods",
        mcp.WithDescription("Lists pods in namespace (uses kubectl)"),
        mcp.WithString("namespace", mcp.Required(), mcp.Description("Target namespace")),
        mcp.WithString("label_selector", mcp.Description("Filter by labels (e.g., app=myapp)")),
        mcp.WithString("output", mcp.Description("Output format: table, json, yaml")),
    )
}

func DevspaceListPodsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    namespace := request.Params.Arguments["namespace"].(string)
    args := []string{"get", "pods", "-n", namespace}

    if labels, ok := request.Params.Arguments["label_selector"].(string); ok && labels != "" {
        args = append(args, "-l", labels)
    }

    output := "table"
    if o, ok := request.Params.Arguments["output"].(string); ok {
        output = o
    }
    args = append(args, "-o", output)

    result := executor.Execute("kubectl", args...)
    return mcp.NewToolResultText(result.Stdout), nil
}
```

**Option B: Parse from devspace analyze**
Not recommended - analyze doesn't expose pod list directly.

---

### 4. No Exec/Shell Capability

**Status:** ✅ Feasible | ✅ IMPLEMENTED

**Problem:** Can't run commands inside pods

**CLI Support:** `devspace enter` supports non-interactive execution with `--tty=false`

```bash
devspace enter --tty=false -- cat /etc/hosts
devspace enter --tty=false --pod=mypod -- curl localhost:8080
```

**Implementation:**

```go
func DevspaceExecTool() mcp.Tool {
    return mcp.NewTool("devspace_exec",
        mcp.WithDescription("Execute a command in a container"),
        mcp.WithString("command", mcp.Required(), mcp.Description("Command to execute")),
        mcp.WithString("working_dir", mcp.Description("DevSpace project directory")),
        mcp.WithString("namespace", mcp.Description("Target namespace")),
        mcp.WithString("pod", mcp.Description("Specific pod name")),
        mcp.WithString("container", mcp.Description("Specific container name")),
        mcp.WithString("label_selector", mcp.Description("Label selector (e.g., app=myapp)")),
        mcp.WithString("image_selector", mcp.Description("Image selector (e.g., nginx:latest)")),
        mcp.WithString("workdir", mcp.Description("Working directory inside container")),
    )
}

func DevspaceExecHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    command := request.Params.Arguments["command"].(string)

    args := []string{"enter", "--tty=false", "--pick=false", "--no-colors"}

    if pod, ok := request.Params.Arguments["pod"].(string); ok && pod != "" {
        args = append(args, "--pod", pod)
    }
    if container, ok := request.Params.Arguments["container"].(string); ok && container != "" {
        args = append(args, "-c", container)
    }
    if labelSelector, ok := request.Params.Arguments["label_selector"].(string); ok && labelSelector != "" {
        args = append(args, "-l", labelSelector)
    }
    if imageSelector, ok := request.Params.Arguments["image_selector"].(string); ok && imageSelector != "" {
        args = append(args, "--image-selector", imageSelector)
    }
    if workdir, ok := request.Params.Arguments["workdir"].(string); ok && workdir != "" {
        args = append(args, "--workdir", workdir)
    }

    // Add the command after --
    args = append(args, "--", command)

    // Use extended timeout for exec commands
    opts := executor.Options{Timeout: 5 * time.Minute}
    if workingDir, ok := request.Params.Arguments["working_dir"].(string); ok {
        opts.WorkingDir = workingDir
    }

    result := executor.ExecuteWithOptions("devspace", opts, args...)
    return mcp.NewToolResultText(result.Stdout), nil
}
```

**Key flags:**
- `--tty=false` - Disable TTY for non-interactive use
- `--pick=false` - Disable interactive pod selection
- `--` separates devspace flags from the command

---

### 5. No devspace_dev Command

**Status:** ❌ Not Feasible

**Problem:** Can't start a devspace dev session through MCP

**Why Not Feasible:**
1. `devspace dev` runs continuously until Ctrl+C
2. Requires terminal for interactive features (log streaming, sync status)
3. No way to gracefully stop without terminal signal
4. MCP tools are request/response, not long-running processes

**Alternative:** Expose `devspace run-pipeline` for non-interactive pipelines:

```go
func DevspaceRunPipelineTool() mcp.Tool {
    return mcp.NewTool("devspace_run_pipeline",
        mcp.WithDescription("Execute a DevSpace pipeline (non-interactive alternative to dev)"),
        mcp.WithString("pipeline", mcp.Required(), mcp.Description("Pipeline name to execute")),
        mcp.WithString("working_dir", mcp.Description("DevSpace project directory")),
        // ... other flags
    )
}
```

---

### 6. No Port-Forward Capability

**Status:** ⚠️ Partially Feasible

**Problem:** Can't expose services locally for debugging

**Analysis:**
- `devspace list ports` - Shows configured port forwards (works)
- No standalone `devspace port-forward` command - port forwarding is managed by `devspace dev`

**Implementation - Part 1: List Ports (Easy)**

```go
func DevspaceListPortsTool() mcp.Tool {
    return mcp.NewTool("devspace_list_ports",
        mcp.WithDescription("Lists configured port forwarding rules"),
        mcp.WithString("working_dir", mcp.Description("DevSpace project directory")),
        mcp.WithString("output", mcp.Description("Output format: table or json")),
    )
}

func DevspaceListPortsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := []string{"list", "ports", "--no-colors"}

    if output, ok := request.Params.Arguments["output"].(string); ok && output == "json" {
        args = append(args, "-o", "json")
    }

    result := executor.ExecuteInDir("devspace", workingDir, args...)
    return mcp.NewToolResultText(result.Stdout), nil
}
```

**Implementation - Part 2: Actual Port Forwarding (kubectl wrapper)**

```go
func KubectlPortForwardTool() mcp.Tool {
    return mcp.NewTool("kubectl_port_forward",
        mcp.WithDescription("Forward local port to a pod (uses kubectl, runs in background)"),
        mcp.WithString("namespace", mcp.Required(), mcp.Description("Target namespace")),
        mcp.WithString("target", mcp.Required(), mcp.Description("Pod name or service/name")),
        mcp.WithString("ports", mcp.Required(), mcp.Description("Port mapping (e.g., 8080:80)")),
    )
}
```

**Note:** Port forwarding is inherently long-running. Would need background process management.

---

### 7. Error Messages Lack Context

**Status:** ✅ Feasible | ✅ IMPLEMENTED

**Problem:** AWS SSO expiration shows as generic "kube config" error

**Implementation:**

```go
// tools/errors.go
package tools

import (
    "strings"
)

type ErrorContext struct {
    Pattern     string
    Suggestion  string
}

var errorPatterns = []ErrorContext{
    {
        Pattern:    "token has expired",
        Suggestion: "Your AWS SSO session has expired. Run: aws sso login",
    },
    {
        Pattern:    "Unable to connect to the server",
        Suggestion: "Cannot reach Kubernetes cluster. Check your VPN connection or cluster status.",
    },
    {
        Pattern:    "forbidden",
        Suggestion: "Permission denied. Verify your RBAC permissions for this namespace.",
    },
    {
        Pattern:    "not found",
        Suggestion: "Resource not found. Verify the namespace and resource name.",
    },
    {
        Pattern:    "context deadline exceeded",
        Suggestion: "Operation timed out. The cluster may be under heavy load or unreachable.",
    },
    {
        Pattern:    "devspace.yaml",
        Suggestion: "No devspace.yaml found. Use working_dir parameter to specify the project location.",
    },
}

func EnhanceError(stderr string) string {
    for _, ec := range errorPatterns {
        if strings.Contains(strings.ToLower(stderr), strings.ToLower(ec.Pattern)) {
            return stderr + "\n\n💡 Suggestion: " + ec.Suggestion
        }
    }
    return stderr
}
```

**Usage in handlers:**

```go
if result.ExitCode != 0 {
    enhancedError := tools.EnhanceError(result.Stderr)
    return mcp.NewToolResultError(enhancedError), nil
}
```

---

### 8. No Streaming/Progress for Long Operations

**Status:** ❌ Not Feasible

**Problem:** `devspace build` and `devspace deploy` can take minutes without feedback

**Why Not Feasible:**
- MCP protocol is request/response based
- No streaming support in current MCP specification
- Would require protocol-level changes

**Workarounds:**

1. **Document expected durations** in tool descriptions
2. **Increase timeouts** appropriately (already done: 10 min for build/deploy)
3. **Return partial output** if command times out

```go
// In tool description
mcp.WithDescription("Builds images and pushes them. Note: This operation may take several minutes for large projects.")
```

---

### 9. Missing Health Check / Status Command

**Status:** ✅ Feasible | ✅ IMPLEMENTED

**Problem:** No quick way to check if devspace environment is healthy

**Implementation:** Composite command aggregating multiple sources

```go
func DevspaceStatusTool() mcp.Tool {
    return mcp.NewTool("devspace_status",
        mcp.WithDescription("Shows overall DevSpace environment health status"),
        mcp.WithString("working_dir", mcp.Required(), mcp.Description("DevSpace project directory")),
        mcp.WithString("namespace", mcp.Description("Target namespace")),
    )
}

func DevspaceStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    workingDir := request.Params.Arguments["working_dir"].(string)
    namespace, _ := request.Params.Arguments["namespace"].(string)

    var status strings.Builder
    status.WriteString("# DevSpace Environment Status\n\n")

    // 1. Check devspace.yaml exists
    configPath := filepath.Join(workingDir, "devspace.yaml")
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        status.WriteString("❌ No devspace.yaml found\n")
        return mcp.NewToolResultText(status.String()), nil
    }
    status.WriteString("✅ devspace.yaml found\n\n")

    // 2. Get deployments status
    status.WriteString("## Deployments\n")
    args := []string{"list", "deployments", "--no-colors"}
    if namespace != "" {
        args = append(args, "-n", namespace)
    }
    result := executor.ExecuteInDir("devspace", workingDir, args...)
    if result.ExitCode == 0 {
        status.WriteString(result.Stdout)
    } else {
        status.WriteString("⚠️ Could not fetch deployments: " + result.Stderr + "\n")
    }

    // 3. Run analyze
    status.WriteString("\n## Analysis\n")
    args = []string{"analyze", "--no-colors", "--timeout=30"}
    if namespace != "" {
        args = append(args, "-n", namespace)
    }
    result = executor.ExecuteInDir("devspace", workingDir, args...)
    if result.ExitCode == 0 {
        if strings.TrimSpace(result.Stdout) == "" {
            status.WriteString("✅ No issues detected\n")
        } else {
            status.WriteString(result.Stdout)
        }
    } else {
        status.WriteString("⚠️ Analysis failed: " + result.Stderr + "\n")
    }

    // 4. List configured sync paths
    status.WriteString("\n## Configured Sync Paths\n")
    result = executor.ExecuteInDir("devspace", workingDir, "list", "sync", "--no-colors")
    status.WriteString(result.Stdout)

    // 5. List configured ports
    status.WriteString("\n## Configured Port Forwards\n")
    result = executor.ExecuteInDir("devspace", workingDir, "list", "ports", "--no-colors")
    status.WriteString(result.Stdout)

    return mcp.NewToolResultText(status.String()), nil
}
```

---

### 10. Log Filtering Options Limited

**Status:** ⚠️ Partially Feasible | ✅ IMPLEMENTED

**Problem:** Can't filter by log level, search within logs, or get multi-container logs

**CLI Capabilities:**
- `--lines` - Number of lines (supported)
- `--container` - Specific container (supported)
- `--label-selector` - Filter pods (supported)
- No grep/level filtering built-in

**Implementation - Enhanced Logs Tool:**

```go
func DevspaceLogsTool() mcp.Tool {
    return mcp.NewTool("devspace_logs",
        mcp.WithDescription("Get container logs with filtering options"),
        mcp.WithString("working_dir", mcp.Description("DevSpace project directory")),
        mcp.WithString("namespace", mcp.Description("Target namespace")),
        mcp.WithString("pod", mcp.Description("Specific pod name")),
        mcp.WithString("container", mcp.Description("Specific container name")),
        mcp.WithString("label_selector", mcp.Description("Label selector (e.g., app=myapp)")),
        mcp.WithNumber("lines", mcp.Description("Number of lines to return (default 200)")),
        mcp.WithString("grep", mcp.Description("Filter logs containing this string")),
        mcp.WithString("grep_level", mcp.Description("Filter by log level: error, warn, info")),
    )
}

func DevspaceLogsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := []string{"logs", "--no-colors", "--pick=false"}

    // Standard devspace flags
    if lines, ok := request.Params.Arguments["lines"].(float64); ok {
        args = append(args, fmt.Sprintf("--lines=%d", int(lines)))
    }
    if pod, ok := request.Params.Arguments["pod"].(string); ok && pod != "" {
        args = append(args, "--pod", pod)
    }
    if container, ok := request.Params.Arguments["container"].(string); ok && container != "" {
        args = append(args, "-c", container)
    }
    if labelSelector, ok := request.Params.Arguments["label_selector"].(string); ok && labelSelector != "" {
        args = append(args, "-l", labelSelector)
    }

    result := executor.ExecuteInDir("devspace", workingDir, args...)
    output := result.Stdout

    // Post-process: grep filter
    if grep, ok := request.Params.Arguments["grep"].(string); ok && grep != "" {
        output = filterLines(output, grep)
    }

    // Post-process: level filter
    if level, ok := request.Params.Arguments["grep_level"].(string); ok && level != "" {
        output = filterByLevel(output, level)
    }

    return mcp.NewToolResultText(output), nil
}

func filterLines(input, pattern string) string {
    var filtered []string
    for _, line := range strings.Split(input, "\n") {
        if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
            filtered = append(filtered, line)
        }
    }
    return strings.Join(filtered, "\n")
}

func filterByLevel(input, level string) string {
    patterns := map[string][]string{
        "error": {"error", "err", "fatal", "panic"},
        "warn":  {"warn", "warning"},
        "info":  {"info"},
    }

    levelPatterns, ok := patterns[strings.ToLower(level)]
    if !ok {
        return input
    }

    var filtered []string
    for _, line := range strings.Split(input, "\n") {
        lineLower := strings.ToLower(line)
        for _, p := range levelPatterns {
            if strings.Contains(lineLower, p) {
                filtered = append(filtered, line)
                break
            }
        }
    }
    return strings.Join(filtered, "\n")
}
```

**Limitations:**
- Multi-container logs require multiple calls (one per container)
- Level filtering is heuristic-based (looks for keywords)

---

## Implementation Priority

### High Priority (High Impact, Low-Medium Effort)

1. **devspace_exec** - Essential for debugging
2. **devspace_list_pods** - Frequently needed
3. **Better error messages** - Improves DX significantly
4. **Consistent working_dir** - Reduces confusion

### Medium Priority

5. **devspace_status** - Useful health overview
6. **Enhanced logs (grep/level)** - Quality of life improvement
7. **devspace_list_ports** - Completes list commands

### Low Priority / Deferred

8. **Better analyze output** - Limited by CLI
9. **Port forwarding** - Complex, needs background process management
10. **devspace_dev** - Not feasible
11. **Streaming progress** - Not feasible with MCP

---

## Files to Create/Modify

| File | Changes |
|------|---------|
| `tools/exec.go` | New file: devspace_exec tool |
| `tools/pods.go` | New file: devspace_list_pods (kubectl wrapper) |
| `tools/status.go` | New file: devspace_status composite tool |
| `tools/ports.go` | New file: devspace_list_ports tool |
| `tools/errors.go` | New file: error enhancement utilities |
| `tools/logs.go` | Modify: add grep/level filtering |
| `tools/analyze.go` | Modify: add patient/timeout flags |
| `tools/validate.go` | Modify: add devspace.yaml check |
| `tools/tools.go` | Modify: register new tools |
