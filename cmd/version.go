package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(ver, date, sha string) *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			if asJSON {
				b, _ := json.MarshalIndent(map[string]string{
					"version":   ver,
					"buildDate": date,
					"commit":    sha,
				}, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Printf("version:    %s\nbuild date: %s\ncommit:     %s\n", ver, date, sha)
			}
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "output as JSON")
	return cmd
}