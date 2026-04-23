package next

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type choice struct {
	label string
	value string
}

func askChoice(question string, choices []choice) string {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n? %s:\n", question)

	for i, c := range choices {
		fmt.Printf("  %d) %s\n", i+1, c.label)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nEnter choice [1-%d]: ", len(choices))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		n, err := strconv.Atoi(input)
		if err == nil && n >= 1 && n <= len(choices) {
			chosen := choices[n-1]
			color.New(color.FgGreen).Printf("  ✓ %s\n", chosen.label)
			return chosen.value
		}
		color.New(color.FgRed).Printf("  Invalid choice, enter 1-%d\n", len(choices))
	}
}

func askBool(question string) bool {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n? %s [y/N]: ", question)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	result := input == "y" || input == "yes"
	if result {
		color.New(color.FgGreen).Println("  ✓ Yes")
	} else {
		color.New(color.FgYellow).Println("  ✗ No")
	}
	return result
}

func Prompt(projectName string) (*NextOptions, error) {
	fmt.Printf("\nConfiguring ")
	color.New(color.FgCyan, color.Bold).Printf("%s\n", projectName)
	fmt.Println(strings.Repeat("─", 40))

	opts := &NextOptions{ProjectName: projectName}
	opts.SetDefaults()

	opts.DB = askChoice("Database", []choice{
		{"Prisma + SQLite  (local, zero config)", "prisma-sqlite"},
		{"Prisma + PostgreSQL", "prisma-pg"},
		{"Drizzle + PostgreSQL", "drizzle-pg"},
		{"Drizzle + MySQL", "drizzle-mysql"},
		{"None", "none"},
	})

	opts.Auth = askChoice("Authentication", []choice{
		{"Better Auth — email + password", "better-auth-email"},
		{"Better Auth — email + password + GitHub OAuth", "better-auth-github"},
		{"None", "none"},
	})

	opts.State = askBool("Include Zustand (state management)?")
	opts.Docker = askBool("Include Docker + docker-compose?")

	return opts, nil
}