package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"promptc/internal/ui"
)

var (
	buildVersion string
	buildCommit  string
	buildDate    string
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ver := buildVersion
			if ver == "" {
				ver = "dev"
			}

			fmt.Println(ui.StyleBold.Render(fmt.Sprintf("promptc %s", ver)))

			if buildCommit != "" {
				fmt.Println(ui.StyleContent.Render(fmt.Sprintf("Commit:  %s", buildCommit)))
			}

			dateStr := formatBuildDate(buildDate)
			if dateStr != "" {
				fmt.Println(ui.StyleContent.Render(fmt.Sprintf("Built:   %s", dateStr)))
			}
		},
	}
}
