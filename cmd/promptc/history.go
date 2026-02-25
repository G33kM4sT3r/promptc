package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"promptc/internal/history"
	"promptc/internal/render"
	"promptc/internal/score"
)

var (
	historyOutputFlag string
	historyLimitFlag  int
	historySearchFlag string
	historyDeleteFlag int
	historyClearFlag  bool
	historyExportFlag bool
	historyYesFlag    bool
)

var historyCmd = &cobra.Command{
	Use:   "history [index]",
	Short: "List or recall prompt history",
	Long:  "List recent prompt compilations or recall a specific entry by index.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVarP(&historyOutputFlag, "output", "o", "text", "Output format (text, json, yaml)")
	historyCmd.Flags().IntVar(&historyLimitFlag, "limit", 0, "Limit number of entries shown (0 = all)")
	historyCmd.Flags().StringVar(&historySearchFlag, "search", "", "Search history by input text")
	historyCmd.Flags().IntVar(&historyDeleteFlag, "delete", -1, "Delete entry at index")
	historyCmd.Flags().BoolVar(&historyClearFlag, "clear", false, "Clear all history")
	historyCmd.Flags().BoolVar(&historyExportFlag, "export", false, "Export full history as JSON")
	historyCmd.Flags().BoolVar(&historyYesFlag, "yes", false, "Skip confirmation prompts")
}

func runHistory(cmd *cobra.Command, args []string) error {
	_, baseDir, err := loadConfig()
	if err != nil {
		return err
	}

	store := history.NewStore(baseDir)
	timeFmt := loadTranslator(baseDir, langFlag)

	// Handle --clear
	if historyClearFlag {
		if !historyYesFlag {
			fmt.Print("Clear all history? [y/N] ")
			var answer string
			_, _ = fmt.Scanln(&answer)
			if !strings.EqualFold(answer, "y") {
				fmt.Println("Canceled.")
				return nil
			}
		}
		if err := store.Clear(); err != nil {
			return err
		}
		fmt.Println("History cleared.")
		return nil
	}

	// Handle --delete
	if historyDeleteFlag >= 0 {
		if !historyYesFlag {
			fmt.Printf("Delete history entry %d? [y/N] ", historyDeleteFlag)
			var answer string
			_, _ = fmt.Scanln(&answer)
			if !strings.EqualFold(answer, "y") {
				fmt.Println("Canceled.")
				return nil
			}
		}
		if err := store.Delete(historyDeleteFlag); err != nil {
			return err
		}
		fmt.Printf("Deleted entry %d.\n", historyDeleteFlag)
		return nil
	}

	// Handle --export
	if historyExportFlag {
		entries, err := store.List()
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	// Handle --search
	if historySearchFlag != "" {
		results := store.Search(historySearchFlag)
		if len(results) == 0 {
			fmt.Println("No matching entries.")
			return nil
		}
		for i, e := range results {
			ts := timeFmt.FormatDateTime(e.Timestamp)
			fmt.Printf("[%d] %s  %s  (score: %d, lang: %s)\n", i, ts, e.Input, e.Score, e.Language)
		}
		return nil
	}

	// List history
	if len(args) == 0 {
		entries, err := store.List()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("No history entries.")
			return nil
		}
		if historyLimitFlag > 0 && len(entries) > historyLimitFlag {
			entries = entries[len(entries)-historyLimitFlag:]
		}
		for i, e := range entries {
			ts := timeFmt.FormatDateTime(e.Timestamp)
			fmt.Printf("[%d] %s  %s  (score: %d, lang: %s)\n", i, ts, e.Input, e.Score, e.Language)
		}
		return nil
	}

	// Recall by index
	index, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid index: %s", args[0])
	}

	entry, err := store.Get(index)
	if err != nil {
		return err
	}

	translator := loadTranslator(baseDir, entry.Language)

	var r render.Renderer
	switch historyOutputFlag {
	case "json":
		r = &render.JSONRenderer{}
	case "yaml":
		r = &render.YAMLRenderer{}
	default:
		r = render.NewTranslated(translator)
	}

	fmt.Println(r.Render(entry.Spec))

	scoreResult := score.Score(entry.Spec)
	sb := render.ScoreBreakdown{
		Total:      scoreResult.Total,
		Breakdown:  scoreResult.Breakdown,
		MaxWeights: score.MaxWeights(),
	}
	switch historyOutputFlag {
	case "json":
		fmt.Fprintln(os.Stderr, (&render.JSONRenderer{}).RenderScore(sb))
	case "yaml":
		fmt.Fprintln(os.Stderr, (&render.YAMLRenderer{}).RenderScore(sb))
	default:
		fmt.Fprint(os.Stderr, render.NewTranslated(translator).RenderScore(sb))
	}
	return nil
}
