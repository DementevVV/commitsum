package github

import (
	"fmt"
	"os/exec"
	"strings"
)

// Error represents a GitHub CLI error.
type Error struct {
	Command string
	Output  string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("GitHub CLI error: %v\nCommand: %s\nOutput: %s", e.Err, e.Command, e.Output)
}

// IsAuthError checks if the error is an authentication issue.
func (e *Error) IsAuthError() bool {
	output := strings.ToLower(e.Output)
	return strings.Contains(output, "authentication") ||
		strings.Contains(output, "not logged in") ||
		strings.Contains(output, "unauthorized")
}

// IsRateLimitError checks if the error is a rate limit issue.
func (e *Error) IsRateLimitError() bool {
	output := strings.ToLower(e.Output)
	return strings.Contains(output, "rate limit") ||
		strings.Contains(output, "api rate limit exceeded")
}

// IsNetworkError checks if the error is a network issue.
func (e *Error) IsNetworkError() bool {
	output := strings.ToLower(e.Output)
	return strings.Contains(output, "network") ||
		strings.Contains(output, "connection") ||
		strings.Contains(output, "timeout") ||
		strings.Contains(output, "dns")
}

// WrapError wraps a GitHub CLI command error.
func WrapError(cmd *exec.Cmd, output []byte, err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		Command: strings.Join(cmd.Args, " "),
		Output:  strings.TrimSpace(string(output)),
		Err:     err,
	}
}

// GetUserFriendlyMessage returns a user-friendly error message.
func GetUserFriendlyMessage(err error) string {
	if ghErr, ok := err.(*Error); ok {
		if ghErr.IsAuthError() {
			return "GitHub authentication required. Run 'gh auth login' to authenticate."
		}
		if ghErr.IsRateLimitError() {
			return "GitHub API rate limit exceeded. Please wait and try again later."
		}
		if ghErr.IsNetworkError() {
			return "Network error. Please check your internet connection and try again."
		}
	}

	errStr := err.Error()
	if strings.Contains(errStr, "executable file not found") {
		return "GitHub CLI (gh) is not installed. Please install it from https://cli.github.com/"
	}

	return fmt.Sprintf("Error: %v", err)
}
