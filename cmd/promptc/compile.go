package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"time"

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
	Args:  cobra.ExactArgs(1),
	RunE:  runCompile,
}

func init() {
	compileCmd.Flags().BoolVarP(&explainFlag, "explain", "e", false, "Explain how the input was interpreted")
	compileCmd.Flags().StringVarP(&outputFlag, "output", "o", "text", "Output format (text, json, yaml)")
	compileCmd.Flags().BoolVar(&scoreFlag, "score", false, "Show quality score for the generated prompt")
	compileCmd.Flags().BoolVar(&copyFlag, "copy", false, "Copy output to clipboard")
}

func runCompile(cmd *cobra.Command, args []string) error {
	input := args[0]

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
			result := pipeline.ApplyRulesWithTrace(s, translator)
			if outputFlag == "text" {
				explain.Print(s, &result)
				fmt.Print("\nDerived Prompt:\n\n")
			}
			spec = result.Spec
			return r.Render(result.Spec), nil
		}

		spec = pipeline.ApplyRules(s, translator)
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
		fmt.Fprintf(os.Stderr, "\nScore: %d/100\n", scoreResult.Total)
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

func loadTranslator(baseDir, lang string) *i18n.Translator {
	if lang == "" || lang == "unknown" {
		lang = "en"
	}
	translationsDir := filepath.Join(baseDir, "translations")
	translator, err := i18n.Load(translationsDir, lang, "en")
	if err != nil {
		translator, _ = i18n.Load(translationsDir, "en", "en")
	}
	return translator
}
