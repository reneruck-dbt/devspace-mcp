# DevSpace MCP Server Feedback

Based on my experience debugging the integration test failures, here are observations on the devspace MCP server:

## What Worked Well

- devspace_list_deployments - Clear output showing deployment status
- devspace_logs - Useful for getting pod logs with label selectors
- devspace_list_profiles - Quick way to see available profiles

## Shortcomings & Missing Features

1. Inconsistent working_dir requirement
- Some tools work without it, others fail with "Cannot find a devspace.yaml"
- Should either auto-detect from cwd or make the requirement consistent

2. devspace_analyze returns minimal output
- Only showed "Checking status..." with no actual analysis
- Expected: pod health, resource usage, warnings, or issues detected

3. No pod listing capability
- Had to fall back to kubectl get pods to see pod status
- Missing: devspace_list_pods or similar

4. No exec/shell capability
- Couldn't run commands inside pods (needed curl to debug GraphQL)
- Missing: devspace_exec or devspace_enter

5. No devspace_dev command
- Can't start a devspace session through MCP
- This is a primary workflow operation

6. No port-forward capability
- Can't expose services locally for debugging
- Missing: devspace_port_forward

7. Error messages lack context
- AWS SSO expiration showed as generic "kube config" error
- Better: detect common auth issues and suggest fixes

8. No streaming/progress for long operations
- devspace_deploy and devspace_build can take minutes
- Would benefit from progress updates or streaming output

9. Missing health check / status command
- No quick way to check if devspace environment is healthy
- Missing: devspace_status showing overall health

10. Log filtering options limited
- Can't filter by log level (error, warn, info)
- Can't search/grep within logs
- Can't get logs from multiple containers easily

Suggested Additions
┌───────────────────────┬──────────────────────────────────────┐
│         Tool          │               Purpose                │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_list_pods    │ List pods with status, restarts, age │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_exec         │ Execute command in a pod/container   │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_status       │ Overall devspace health check        │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_dev          │ Start devspace dev session           │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_port_forward │ Forward local port to service        │
├───────────────────────┼──────────────────────────────────────┤
│ devspace_sync_status  │ Check file sync status               │
└───────────────────────┴──────────────────────────────────────┘
