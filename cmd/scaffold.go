package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zorojuro75/devtool/internal/scaffold"
)

func newScaffoldCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "scaffold [framework] [project-name]",
		Short: "Scaffold a new project from a template",
		Example: `  devtool scaffold go myapi
		devtool scaffold laravel blog`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			framework := args[0]
			name := framework
			if len(args) == 2 {
				name = args[1]
			}

			files, err := scaffold.Generate(framework, name, force)
			if err != nil {
				return err
			}

			color.New(color.FgGreen, color.Bold).Printf("\nScaffolded %s project: %s\n\n", framework, name)
			for _, f := range files {
				fmt.Printf("  + %s\n", f)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing directory without prompting")
	return cmd
}