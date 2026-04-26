package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "A developer CLI for scaffolding Next.js fullstack projects",
	Long:  `devtool scaffolds production-ready Next.js 16 fullstack projects with Better Auth, Prisma or Drizzle, Tailwind CSS, and shadcn/ui.`,
}

func Execute(ver, date, sha string) {
	rootCmd.AddCommand(
		newNextCmd(),
		newVersionCmd(ver, date, sha),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colour output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose logging")
}

var (
	noColor bool
	verbose bool
)