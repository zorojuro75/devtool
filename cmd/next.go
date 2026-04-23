package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zorojuro75/devtool/internal/next"
)

func newNextCmd() *cobra.Command {
	var db, auth string
	var state, docker, noPrompt bool

	cmd := &cobra.Command{
		Use:   "next [project-name]",
		Short: "Scaffold a Next.js 16 fullstack project",
		Long: `Scaffold a production-ready Next.js 16 fullstack project with:
  - App Router + TypeScript
  - Tailwind CSS + shadcn/ui
  - Better Auth (email/password or + GitHub OAuth)
  - Prisma or Drizzle ORM
  - Zustand state management (optional)
  - Docker + docker-compose (optional)`,
		Example: `  devtool next myapp
  devtool next myapp --db prisma-sqlite --auth better-auth-email --state --docker --no-prompt`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts *next.NextOptions
			var err error

			if noPrompt {
				if db == "" || auth == "" {
					return fmt.Errorf("--no-prompt requires --db and --auth flags")
				}
				opts = &next.NextOptions{
					ProjectName: args[0],
					DB:          db,
					Auth:        auth,
					State:       state,
					Docker:      docker,
				}
				opts.SetDefaults()
			} else {
				opts, err = next.Prompt(args[0])
				if err != nil {
					return err
				}
				// flags override prompts if provided
				if db != "" {
					opts.DB = db
				}
				if auth != "" {
					opts.Auth = auth
				}
				if state {
					opts.State = true
				}
				if docker {
					opts.Docker = true
				}
			}

			_, err = next.Generate(opts)
			return err
		},
	}

	cmd.Flags().StringVar(&db, "db", "", "prisma-sqlite | prisma-pg | drizzle-pg | drizzle-mysql | none")
	cmd.Flags().StringVar(&auth, "auth", "", "better-auth-email | better-auth-github | none")
	cmd.Flags().BoolVar(&state, "state", false, "include Zustand state management")
	cmd.Flags().BoolVar(&docker, "docker", false, "include Docker + docker-compose")
	cmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "skip interactive prompts, use flags only")
	return cmd
}