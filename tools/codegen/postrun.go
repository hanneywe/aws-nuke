package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// CommandExecutor abstracts command execution for testability.
type CommandExecutor interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// DefaultExecutor runs commands using os/exec.
type DefaultExecutor struct{}

// Run executes a command and returns its combined output.
func (executor *DefaultExecutor) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.CombinedOutput()
}

// RunPostSteps executes goimports, go build, and golangci-lint in sequence.
// It stops on the first failure and returns an error with the tool output.
func RunPostSteps(ctx context.Context, executor CommandExecutor) error {
	steps := []struct {
		name string
		args []string
	}{
		{"goimports", []string{"-w", "./resources/"}},
		{"go", []string{"build", "./resources/..."}},
		{"golangci-lint", []string{"run", "./resources/..."}},
	}

	for _, step := range steps {
		output, err := executor.Run(ctx, step.name, step.args...)
		if err != nil {
			return fmt.Errorf("%s %s failed:\n%s",
				step.name, strings.Join(step.args, " "), string(output))
		}
	}

	return nil
}

// CommitFiles stages the given files and creates a git commit with a
// conventional commit message.
func CommitFiles(ctx context.Context, executor CommandExecutor, files []string, serviceName, resourceName string) error {
	addArgs := append([]string{"add"}, files...)

	output, err := executor.Run(ctx, "git", addArgs...)
	if err != nil {
		return fmt.Errorf("git add failed:\n%s", string(output))
	}

	commitMessage := fmt.Sprintf("feat(%s): add %s resource", serviceName, resourceName)

	output, err = executor.Run(ctx, "git", "commit", "-m", commitMessage)
	if err != nil {
		return fmt.Errorf("git commit failed:\n%s", string(output))
	}

	return nil
}
