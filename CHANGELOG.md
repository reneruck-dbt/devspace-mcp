# Changelog

All notable changes to the DevSpace MCP Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### New Tools

- **devspace_exec** - Execute commands inside containers non-interactively
  - Uses `devspace enter --tty=false` for command execution
  - Supports pod, container, label selector, image selector, and workdir parameters
  - 5-minute timeout for operations
  - Enables debugging commands like `ls`, `curl`, `cat` without terminal access

- **devspace_list_pods** - List pods in Kubernetes namespaces
  - kubectl wrapper implementation (DevSpace doesn't provide this natively)
  - Supports label selectors, field selectors, and multiple output formats (wide, json, yaml, name)
  - All-namespaces flag support
  - Essential for pod discovery and inspection

- **devspace_status** - Comprehensive environment health overview
  - Aggregates information from multiple sources:
    - devspace.yaml validation
    - Deployment status
    - Namespace analysis
    - Configured sync paths
    - Configured port forwards
  - Graceful degradation - continues even if individual checks fail
  - Clear section-based output with status indicators (✅/❌/⚠️)

- **devspace_list_ports** - List configured port forwarding rules
  - Shows port forwarding configuration from devspace.yaml
  - Supports table and JSON output formats
  - Useful for understanding local-to-container port mappings

#### Enhanced Tools

- **devspace_logs** - Added client-side filtering capabilities
  - `grep` parameter for text-based filtering (case-insensitive)
  - `grep_level` parameter for log level filtering (error, warn, info)
  - Error level matches: error, err, fatal, panic, failed, failure
  - Warn level matches: warn, warning
  - Info level matches: info
  - Post-processing after retrieval from devspace CLI

- **devspace_analyze** - Added analysis control flags
  - `patient` flag: wait for all resources to be ready before reporting
  - `ignore_pod_restarts` flag: ignore restart events of running pods
  - Updated description to clarify minimal output when healthy is expected behavior
  - Provides more control over what constitutes a reportable problem

#### Error Handling & Validation

- **Contextual error messages** - Intelligent error pattern detection
  - 17 common error patterns with actionable suggestions:
    - AWS SSO token expiration → suggests `aws sso login`
    - Connection refused → suggests checking VPN/cluster
    - Forbidden errors → suggests checking RBAC permissions
    - Resource not found → suggests verifying namespace/name
    - ImagePullBackOff/CrashLoopBackOff/OOMKilled → container-specific guidance
    - Certificate errors → suggests renewal/verification
    - devspace.yaml not found → suggests using working_dir parameter
  - Applied to all major tool handlers (build, deploy, analyze, logs, exec, pods)
  - Case-insensitive pattern matching
  - Significantly improves developer experience during failures

- **Working directory validation** - Consistent devspace.yaml checking
  - New `ValidateDevspaceYaml()` function in validate.go
  - Checks for devspace.yaml in specified directory or current directory
  - Provides clear error messages with working_dir parameter suggestion
  - Available for reuse across tools that require devspace.yaml

### Changed

- Updated feasibility analysis document to mark implemented features
- Improved tool descriptions to clarify behavior and limitations

### Technical Details

- **Test Coverage**: Added 9 new test files with comprehensive unit tests
- **Code Quality**: All tests passing, `make check` successful
- **Documentation**: Each feature documented with inline comments and tool descriptions
- **Commits**: 8 clean, atomic commits (one per feature) with detailed messages

## [0.1.0] - 2024-01-21

### Added

- Initial DevSpace MCP server implementation
- Core tools:
  - devspace_version - Get DevSpace CLI version
  - devspace_list_namespaces - List Kubernetes namespaces
  - devspace_list_contexts - List Kubernetes contexts
  - devspace_list_deployments - List DevSpace deployments
  - devspace_list_profiles - List DevSpace profiles
  - devspace_list_vars - List DevSpace variables
  - devspace_print - Print rendered configuration
  - devspace_analyze - Analyze namespace for problems
  - devspace_logs - Get container logs
  - devspace_build - Build container images
  - devspace_deploy - Deploy to Kubernetes
  - devspace_purge - Remove deployments
  - devspace_run - Run commands/pipelines
- Executor package for command execution with timeouts
- Input validation for security (flag injection prevention)
- MCP protocol integration via mcp-go SDK
- Comprehensive test suite
- Documentation (README, CLAUDE.md, feasibility analysis)

[Unreleased]: https://github.com/yourusername/devspace-mcp/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yourusername/devspace-mcp/releases/tag/v0.1.0
