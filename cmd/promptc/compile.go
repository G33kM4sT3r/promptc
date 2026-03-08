package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"promptc/internal/clipboard"
	"promptc/internal/explain"
	"promptc/internal/history"
	"promptc/internal/i18n"
	"promptc/internal/pipeline"
	"promptc/internal/prompt"
	"promptc/internal/render"
	"promptc/internal/score"
	"promptc/internal/ui"
)

var (
	explainFlag bool
	outputFlag  string
	scoreFlag   bool
	copyFlag    bool
)

var compileCmd = &cobra.Command{
	Use:   "compile [input]",
	Short: "Compile natural language into a structured prompt",
	Long:  "Compile takes a natural language input string and transforms it into a structured AI prompt using rule-based NLP.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCompile,
}

func init() {
	compileCmd.Flags().BoolVarP(&explainFlag, "explain", "e", false, "Explain how the input was interpreted")
	compileCmd.Flags().StringVarP(&outputFlag, "output", "o", "text", "Output format (text, json, yaml)")
	compileCmd.Flags().BoolVar(&scoreFlag, "score", false, "Show quality score for the generated prompt")
	compileCmd.Flags().BoolVar(&copyFlag, "copy", false, "Copy output to clipboard")
}

func runCompile(cmd *cobra.Command, args []string) error {
	input, err := resolveInput(args)
	if err != nil {
		return err
	}

	cfg, baseDir, err := loadConfig()
	if err != nil {
		return err
	}

	detector, _ := createDetector(baseDir, cfg)
	defer func() { _ = detector.Close() }()

	var spec prompt.PromptSpec
	var detectedLang string

	work := func() (string, error) {
		s, err := pipeline.Run(input, cfg, detector)
		if err != nil {
			return "", err
		}

		detectedLang = s.Language
		if langFlag != "" {
			detectedLang = langFlag
		}
		translator := loadTranslator(baseDir, detectedLang)

		var r render.Renderer
		switch outputFlag {
		case "json":
			r = &render.JSONRenderer{}
		case "yaml":
			r = &render.YAMLRenderer{}
		default:
			r = render.NewTranslated(translator)
		}

		if explainFlag {
			result := pipeline.ApplyRulesWithTrace(s, translator, cfg.Enrichments)
			if outputFlag == "text" {
				explain.Print(s, &result)
				fmt.Print("\nDerived Prompt:\n\n")
			}
			spec = result.Spec
			return r.Render(result.Spec), nil
		}

		spec = pipeline.ApplyRules(s, translator, cfg.Enrichments)
		return r.Render(spec), nil
	}

	output, err := ui.RunWithSpinner("Compiling prompt...", work)
	if err != nil {
		return err
	}

	fmt.Println(output)

	if copyFlag {
		if err := clipboard.Copy(output); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: clipboard copy failed: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "Copied to clipboard.")
		}
	}

	scoreResult := score.Score(spec)

	if scoreFlag {
		sb := render.ScoreBreakdown{
			Total:      scoreResult.Total,
			Breakdown:  scoreResult.Breakdown,
			MaxWeights: score.MaxWeights(),
		}
		switch outputFlag {
		case "json":
			fmt.Fprintln(os.Stderr, (&render.JSONRenderer{}).RenderScore(sb))
		case "yaml":
			fmt.Fprintln(os.Stderr, (&render.YAMLRenderer{}).RenderScore(sb))
		default:
			translator := loadTranslator(baseDir, detectedLang)
			fmt.Fprint(os.Stderr, render.NewTranslated(translator).RenderScore(sb))
		}
	}

	// Auto-save to history
	store := history.NewStore(baseDir)
	_ = store.Add(history.Entry{
		Input:     input,
		Spec:      spec,
		Score:     scoreResult.Total,
		Language:  detectedLang,
		Timestamp: time.Now(),
	})

	return nil
}

func resolveInput(args []string) (string, error) {
	if len(args) == 1 && args[0] != "-" {
		return args[0], nil
	}

	// Read from stdin if piped or dash argument
	if len(args) == 0 || args[0] == "-" {
		if len(args) == 0 && term.IsTerminal(int(os.Stdin.Fd())) {
			return "", fmt.Errorf("no input provided (use an argument or pipe input via stdin)")
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		input := strings.TrimSpace(string(data))
		if input == "" {
			return "", fmt.Errorf("empty input from stdin")
		}
		return input, nil
	}

	return "", fmt.Errorf("no input provided")
}

func loadTranslator(baseDir, lang string) *i18n.Translator {
	if lang == "" || lang == "unknown" {
		lang = "en"
	}
	langDir := filepath.Join(baseDir, "languages")
	translator, err := i18n.Load(langDir, lang, "en")
	if err != nil {
		translator, _ = i18n.Load(langDir, "en", "en")
	}
	return translator
}
