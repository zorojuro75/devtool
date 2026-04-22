package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Log(since string, maxCommits int, branch string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := []string{"log", "--oneline", "--no-merges", "--since=" + since}
	if branch != "" {
		args = append(args, branch)
	}

	out, err := exec.CommandContext(ctx, "git", args...).Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git log timed out after 10s")
		}
		return "", fmt.Errorf("git error: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > maxCommits {
		lines = lines[:maxCommits]
	}

	return strings.Join(lines, "\n"), nil
}