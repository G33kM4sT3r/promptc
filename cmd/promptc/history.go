package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"promptc/internal/history"
	"promptc/internal/render"
)

var historyOutputFlag string

var historyCmd = &cobra.Command{
	Use:   "history [index]",
	Short: "List or recall prompt history",
	Long:  "List recent prompt compilations or recall a specific entry by index.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVarP(&historyOutputFlag, "output", "o", "text", "Output format (text, json, yaml)")
}

func runHistory(cmd *cobra.Command, args []string) error {
	_, baseDir, err := loadConfig()
	if err != nil {
		return err
	}

	store := history.NewStore(baseDir)

	if len(args) == 0 {
		entries, err := store.List()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("No history entries.")
			return nil
		}
		for i, e := range entries {
			ts := e.Timestamp.Format("2006-01-02 15:04")
			fmt.Printf("[%d] %s  %s  (score: %d, lang: %s)\n", i, ts, e.Input, e.Score, e.Language)
		}
		return nil
	}

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
	fmt.Fprintf(os.Stderr, "\nScore: %d/100\n", entry.Score)
	return nil
}
