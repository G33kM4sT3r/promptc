package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"promptc/internal/config"
	"promptc/internal/language"
	"promptc/internal/ui"
)

var langFlag string

var rootCmd = &cobra.Command{
	Use:   "promptc",
	Short: "Deterministic prompt compiler",
	Long:  "promptc transforms natural language input into structured AI prompts using rule-based NLP.",
}

// Execute configures the root command with build metadata and runs the CLI.
func Execute(version, commit, date string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Run = runWelcome
	rootCmd.PersistentFlags().StringVarP(&langFlag, "lang", "l", "", "Override detected language (e.g., en, de)")
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(replCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.CompletionOptions.HiddenDefaultCmd = false
}

func runWelcome(cmd *cobra.Command, args []string) {
	banner := ` ██████╗ ██████╗  ██████╗ ███╗   ███╗██████╗ ████████╗ ██████╗
 ██╔══██╗██╔══██╗██╔═══██╗████╗ ████║██╔══██╗╚══██╔══╝██╔════╝
 ██████╔╝██████╔╝██║   ██║██╔████╔██║██████╔╝   ██║   ██║
 ██╔═══╝ ██╔══██╗██║   ██║██║╚██╔╝██║██╔═══╝    ██║   ██║
 ██║     ██║  ██║╚██████╔╝██║ ╚═╝ ██║██║        ██║   ╚██████╗
 ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝╚═╝        ╚═╝    ╚═════╝`

	fmt.Println(ui.GradientText(banner, "#b429f9", "#26c6da"))
	fmt.Println()

	ver := rootCmd.Version
	if ver == "" {
		ver = "dev"
	}
	fmt.Println(ui.StyleDim.Render(fmt.Sprintf("  v%s — Deterministic prompt compiler", ver)))
	fmt.Println()
	fmt.Println(ui.StyleBold.Render("  Usage:"))
	fmt.Println(ui.StyleContent.Render(`    promptc compile "explain closures in Go"    Compile a prompt`))
	fmt.Println(ui.StyleContent.Render(`    promptc compile --explain "input"            With rule tracing`))
	fmt.Println(ui.StyleContent.Render(`    promptc repl                                 Interactive mode`))
	fmt.Println(ui.StyleContent.Render(`    promptc version                              Version info`))
	fmt.Println(ui.StyleContent.Render(`    promptc help                                 Full help`))
	fmt.Println(ui.StyleContent.Render(`    promptc completion bash                       Shell completions`))
	fmt.Println()
}

// --- Shared bootstrap helpers ---

func loadConfig() (*config.Config, string, error) {
	baseDir, err := config.FindBaseDir()
	if err != nil {
		return nil, "", fmt.Errorf("finding base dir: %w", err)
	}
	cfg, err := config.Load(baseDir)
	if err != nil {
		return nil, "", fmt.Errorf("loading config: %w", err)
	}
	return cfg, baseDir, nil
}

func createDetector(baseDir string, cfg *config.Config) (*language.Detector, error) {
	modelPath := filepath.Join(baseDir, "data", "lid.176.ftz")
	detector, err := language.NewDetector(modelPath, cfg)
	return detector, err
}

func formatBuildDate(isoDate string) string {
	if isoDate == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", isoDate)
		if err != nil {
			return isoDate
		}
	}
	_, baseDir, cfgErr := loadConfig()
	if cfgErr != nil {
		return t.Format("Jan 2, 2006")
	}
	return loadTranslator(baseDir, langFlag).FormatDate(t)
}
