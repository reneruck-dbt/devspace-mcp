package tools

import (
	"strings"

	"devspace-mcp/executor"
)

// ErrorContext represents a known error pattern and its helpful suggestion
type ErrorContext struct {
	Pattern    string
	Suggestion string
}

// errorPatterns contains known error patterns and helpful suggestions
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
		Pattern:    "connection refused",
		Suggestion: "Cannot connect to Kubernetes API server. Verify the cluster is running and accessible.",
	},
	{
		Pattern:    "forbidden",
		Suggestion: "Permission denied. Verify your RBAC permissions for this namespace/resource.",
	},
	{
		Pattern:    "unauthorized",
		Suggestion: "Authentication failed. Check your kubeconfig credentials and context.",
	},
	{
		Pattern:    "not found",
		Suggestion: "Resource not found. Verify the namespace, resource name, and that the resource exists.",
	},
	{
		Pattern:    "context deadline exceeded",
		Suggestion: "Operation timed out. The cluster may be under heavy load or unreachable.",
	},
	{
		Pattern:    "no such host",
		Suggestion: "DNS resolution failed. Check your network connection and cluster endpoint.",
	},
	{
		Pattern:    "devspace.yaml",
		Suggestion: "No devspace.yaml found. Use working_dir parameter to specify the project location.",
	},
	{
		Pattern:    "Cannot find a devspace.yaml",
		Suggestion: "No devspace.yaml found in current directory. Use working_dir parameter to specify the project location.",
	},
	{
		Pattern:    "no space left on device",
		Suggestion: "Disk space exhausted. Free up disk space or clean up old images/containers.",
	},
	{
		Pattern:    "ImagePullBackOff",
		Suggestion: "Cannot pull container image. Check image name, registry credentials, and network connectivity.",
	},
	{
		Pattern:    "CrashLoopBackOff",
		Suggestion: "Container is repeatedly crashing. Check pod logs for application errors.",
	},
	{
		Pattern:    "OOMKilled",
		Suggestion: "Container was killed due to out of memory. Increase memory limits or optimize application memory usage.",
	},
	{
		Pattern:    "ErrImagePull",
		Suggestion: "Failed to pull container image. Verify image exists and registry is accessible.",
	},
	{
		Pattern:    "namespace does not exist",
		Suggestion: "The specified namespace doesn't exist. Create it first or check the namespace name.",
	},
	{
		Pattern:    "certificate has expired",
		Suggestion: "TLS certificate expired. Renew the certificate or update your kubeconfig.",
	},
	{
		Pattern:    "x509: certificate",
		Suggestion: "Certificate validation error. Check your kubeconfig certificates or use --insecure flag if appropriate.",
	},
}

// EnhanceError analyzes error output and adds helpful suggestions
func EnhanceError(result executor.Result) string {
	// Combine stderr and error for analysis
	errorText := result.Stderr
	if result.Error != "" {
		if errorText != "" {
			errorText += "\n"
		}
		errorText += result.Error
	}

	// If there's no error text, return formatted output
	if errorText == "" {
		return result.FormatOutput()
	}

	// Check for matching patterns
	for _, ec := range errorPatterns {
		if containsIgnoreCase(errorText, ec.Pattern) {
			// Add suggestion to the output
			enhanced := result.FormatOutput()
			if enhanced != "" {
				enhanced += "\n\n"
			}
			enhanced += "💡 Suggestion: " + ec.Suggestion
			return enhanced
		}
	}

	// No pattern matched, return original output
	return result.FormatOutput()
}

// containsIgnoreCase checks if text contains pattern (case-insensitive)
func containsIgnoreCase(text, pattern string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(pattern))
}
