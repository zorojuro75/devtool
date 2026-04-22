package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupRepo creates a temporary git repo with some commits for testing
func setupRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v failed: %v", args, err)
		}
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")

	// Create a file and commit
	f := filepath.Join(dir, "README.md")
	os.WriteFile(f, []byte("hello"), 0644)
	run("add", ".")
	run("commit", "-m", "initial commit")

	return dir
}

func TestLog_ValidRepo(t *testing.T) {
	dir := setupRepo(t)
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	log, err := Log("1d", 50, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log == "" {
		t.Error("expected non-empty log output")
	}
	t.Logf("log output: %s", log)
}

func TestLog_MaxCommits(t *testing.T) {
	dir := setupRepo(t)
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Add more commits
	run := func(msg string) {
		f := filepath.Join(dir, "file.txt")
		os.WriteFile(f, []byte(msg), 0644)
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = dir
		cmd.Run()
		cmd = exec.Command("git", "commit", "-m", msg)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@test.com",
		)
		cmd.Run()
	}

	for i := 0; i < 5; i++ {
		run("commit number extra")
	}

	log, err := Log("1d", 3, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(log), "\n")
	if len(lines) > 3 {
		t.Errorf("got %d lines, want max 3", len(lines))
	}
}