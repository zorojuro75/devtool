package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zorojuro75/devtool/internal/ai"
	"github.com/zorojuro75/devtool/internal/config"
	devgit "github.com/zorojuro75/devtool/internal/git"
)

func newGitlogCmd() *cobra.Command {
	var since string
	var format string
	var branch string
	var maxCommits int

	cmd := &cobra.Command{
		Use:   "gitlog",
		Short: "Summarise recent git commits using AI",
		Example: `  devtool gitlog --since 7d
  devtool gitlog --since 1d --format standup`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			if cfg.APIKey == "" {
				return fmt.Errorf("no API key found — set api_key in ~/.devtool.yaml or DEVTOOL_API_KEY env var")
			}

			log, err := devgit.Log(since, maxCommits, branch)
			if err != nil {
				return fmt.Errorf("git error: %w\n\nMake sure you are inside a git repository and git is installed.", err)
			}
			if log == "" {
				fmt.Println("No commits found in the given time range.")
				return nil
			}

			prompts := map[string]string{
				"summary":   "Summarise these git commits in 3-5 sentences of clear prose, grouped by theme.",
				"changelog": "Convert these commits into a clean changelog with bullet points grouped under Added, Changed, and Fixed headings.",
				"standup":   "Convert these commits into a concise daily standup update in first person (e.g. 'I worked on...'). Keep it to 3-4 sentences.",
			}

			style, ok := prompts[format]
			if !ok {
				return fmt.Errorf("unknown format %q — choose: summary, changelog, standup", format)
			}

			prompt := fmt.Sprintf("%s\n\nCommits:\n%s", style, log)
			system := "You are a helpful engineering assistant. Be concise and professional."

			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(cfg.Timeout)*time.Second)
			defer cancel()

			client := ai.NewOpenRouter(cfg.APIKey, cfg.Model, cfg.Timeout)

			color.New(color.FgCyan).Fprint(os.Stderr, "Summarising commits... ")
			reader, err := client.Stream(ctx, prompt, system)
			fmt.Fprint(os.Stderr, "\r                       \r")
			if err != nil {
				return fmt.Errorf("AI request failed: %w", err)
			}

			color.New(color.FgGreen, color.Bold).Printf("Git summary (%s):\n", format)
			return ai.PrintStream(reader, os.Stdout)
		},
	}

	cmd.Flags().StringVar(&since, "since", "1d", "time range: 1d, 7d, 30d, or YYYY-MM-DD")
	cmd.Flags().StringVar(&format, "format", "summary", "output style: summary, changelog, standup")
	cmd.Flags().StringVar(&branch, "branch", "", "branch to summarise (default: current)")
	cmd.Flags().IntVar(&maxCommits, "max-commits", 50, "max commits to include")
	return cmd
}