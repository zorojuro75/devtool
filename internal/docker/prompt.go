package docker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// ConfirmDetected asks the user to confirm the auto-detected stack
func ConfirmDetected(info *StackInfo) bool {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n? Confirm stack: %s [Y/n]: ", info.DisplayName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "" || input == "y" || input == "yes"
}

// ConfirmORM asks the user to confirm the auto-detected ORM
func ConfirmORM(info *ORMInfo) bool {
	cyan := color.New(color.FgCyan, color.Bold)
	dim := color.New(color.FgHiBlack)

	if info.Name == "none" {
		dim.Println("\n  No ORM detected in package.json")
		return true
	}

	cyan.Printf("\n? Detected ORM: %s. Confirm? [Y/n]: ", info.DisplayName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "" || input == "y" || input == "yes"
}

// AskServices shows a numbered multi-select prompt for services
func AskServices() ([]string, error) {
	cyan := color.New(color.FgCyan, color.Bold)
	dim := color.New(color.FgHiBlack)
	green := color.New(color.FgGreen)

	cyan.Println("\n? Select services to include:")

	for i, name := range ServiceOrder {
		displayName := ServiceDisplayNames[name]
		fmt.Printf("  %d) %s\n", i+1, displayName)
	}
	fmt.Printf("  %d) None\n", len(ServiceOrder)+1)

	dim.Print("\nEnter choices (comma separated, e.g. 1,3) or press Enter for none: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	noneIndex := len(ServiceOrder) + 1
	if input == "" || input == strconv.Itoa(noneIndex) {
		dim.Println("  No services selected")
		return []string{}, nil
	}

	parts := strings.Split(input, ",")
	selected := []string{}
	seen := map[string]bool{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 || n > noneIndex {
			return nil, fmt.Errorf("invalid choice %q - enter numbers between 1 and %d", part, noneIndex)
		}
		if n == noneIndex {
			return []string{}, nil
		}
		serviceName := ServiceOrder[n-1]
		if !seen[serviceName] {
			selected = append(selected, serviceName)
			seen[serviceName] = true
		}
	}

	fmt.Println()
	for _, s := range selected {
		green.Printf("  + %s\n", ServiceDisplayNames[s])
	}

	return selected, nil
}

// AskPort asks the user for the app port with a default
// Re-prompts on invalid input instead of silently falling back
func AskPort(defaultPort int) int {
	cyan := color.New(color.FgCyan, color.Bold)
	red := color.New(color.FgRed)

	reader := bufio.NewReader(os.Stdin)

	for {
		cyan.Printf("\n? App port [%d]: ", defaultPort)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Empty input = use default
		if input == "" {
			return defaultPort
		}

		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > 65535 {
			red.Printf("  Invalid port %q - enter a number between 1 and 65535\n", input)
			continue
		}

		return n
	}
}

// AskOverwrite asks the user whether to overwrite an existing file
func AskOverwrite(filename string) bool {
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("\n? %s already exists. Overwrite? [y/N]: ", filename)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "y" || input == "yes"
}