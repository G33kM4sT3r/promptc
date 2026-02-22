package main

import (
	"github.com/spf13/cobra"

	"promptc/internal/repl"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start interactive REPL mode",
	Long:  "Start an interactive prompt compilation loop with colored output and REPL commands.",
	Args:  cobra.NoArgs,
	RunE:  runRepl,
}

func runRepl(cmd *cobra.Command, args []string) error {
	cfg, baseDir, err := loadConfig()
	if err != nil {
		return err
	}

	detector, _ := createDetector(baseDir, cfg)
	defer func() { _ = detector.Close() }()

	return repl.Run(cfg, detector, baseDir, langFlag, rootCmd.Version)
}
