package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	noColor  bool
	verbose  bool
	provider string
)

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "A developer CLI with AI assistance",
	Long: `devtool helps you scaffold projects, explain errors, and summarise
git history — all with AI assistance powered by OpenRouter.`,
}

func Execute(ver, date, sha string) {
	rootCmd.AddCommand(
		newScaffoldCmd(),
		newExplainCmd(),
		newGitlogCmd(),
		newVersionCmd(ver, date, sha),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.devtool.yaml)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colour output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose logging")
}