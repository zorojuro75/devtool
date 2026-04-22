package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zorojuro75/devtool/internal/ai"
	"github.com/zorojuro75/devtool/internal/config"
)

func newExplainCmd() *cobra.Command {
	var lang string

	cmd := &cobra.Command{
		Use:   "explain [error message]",
		Short: "Explain an error message using AI",
		Example: `  devtool explain "panic: runtime error: index out of range"
  cat error.log | devtool explain --lang go`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}
			if cfg.APIKey == "" {
				return fmt.Errorf("no API key found — set api_key in ~/.devtool.yaml or DEVTOOL_API_KEY env var")
			}

			// Accept input from arg or stdin
			var input string
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				b, _ := io.ReadAll(os.Stdin)
				input = string(b)
			} else if len(args) > 0 {
				input = strings.Join(args, " ")
			} else {
				return fmt.Errorf("provide an error message as an argument or pipe it via stdin")
			}

			system := `You are a senior software engineer. When given an error message:
1. Identify the error type and language
2. Explain the root cause in plain English (2-3 sentences)
3. Suggest 2-3 concrete fixes with brief code examples
Keep your response concise and practical.`

			if lang != "" {
				system += fmt.Sprintf("\nThe error is from a %s project.", lang)
			}

			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(cfg.Timeout)*time.Second)
			defer cancel()

			client := ai.NewOpenRouter(cfg.APIKey, cfg.Model, cfg.Timeout)

			color.New(color.FgCyan).Fprint(os.Stderr, "Thinking... ")
			reader, err := client.Stream(ctx, input, system)
			fmt.Fprint(os.Stderr, "\r              \r") // clear spinner
			if err != nil {
				return fmt.Errorf("AI request failed: %w", err)
			}

			color.New(color.FgGreen, color.Bold).Println("Explanation:")
			return ai.PrintStream(reader, os.Stdout)
		},
	}

	cmd.Flags().StringVar(&lang, "lang", "", "hint the programming language (e.g. go, php, typescript)")
	return cmd
}