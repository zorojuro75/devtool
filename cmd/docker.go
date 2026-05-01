package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/zorojuro75/devtool/internal/docker"
)

func newDockerCmd() *cobra.Command {
    dockerCmd := &cobra.Command{
        Use:   "docker",
        Short: "Add Docker to your project",
        Long:  `Generate a production-ready Dockerfile and docker-compose.yml for your project.`,
    }

    dockerCmd.AddCommand(newDockerInitCmd())
    return dockerCmd
}

func newDockerInitCmd() *cobra.Command {
    var stack string
    var port int
    var services []string
    var noPrompt bool

    cmd := &cobra.Command{
        Use:   "init",
        Short: "Generate Dockerfile and docker-compose.yml for the current project",
        Example: `  devtool docker init
  devtool docker init --services postgres,redis --no-prompt`,
        RunE: func(cmd *cobra.Command, args []string) error {
            dir, err := os.Getwd()
            if err != nil {
                return fmt.Errorf("cannot get current directory: %w", err)
            }

            opts := docker.NewDockerOptions()

            if stack != "" {
                opts.Stack = stack
                opts.ProjectName = projectNameFromDir(dir)
            } else {
                info, err := docker.DetectOrConfirm(dir)
                if err != nil {
                    return err
                }

                fmt.Printf("Detected: %s\n", info.DisplayName)

                if !noPrompt {
                    if !docker.ConfirmDetected(info) {
                        return fmt.Errorf("aborted")
                    }
                }

                opts.Stack = info.Name
                opts.Port = info.DefaultPort
                opts.ProjectName = projectNameFromDir(dir)
            }

            if noPrompt {
                opts.Services = services
            } else {
                selected, err := docker.AskServices()
                if err != nil {
                    return err
                }
                opts.Services = selected
            }

            if port != 0 {
                opts.Port = port
            } else if !noPrompt {
                opts.Port = docker.AskPort(opts.Port)
            }

            return docker.Generate(opts, dir)
        },
    }

    cmd.Flags().StringVar(&stack, "stack", "", "project stack: nextjs (skip auto-detection)")
    cmd.Flags().IntVar(&port, "port", 0, "app port (default: 3000)")
    cmd.Flags().StringSliceVar(&services, "services", nil, "services: postgres,mysql,redis,mailhog")
    cmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "skip prompts, use flags only")

    return cmd
}

func projectNameFromDir(dir string) string {
    parts := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
    if len(parts) == 0 {
        return "myapp"
    }
    return parts[len(parts)-1]
}
