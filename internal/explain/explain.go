package explain

import (
	"fmt"
	"strings"

	"promptc/internal/rules"
	"promptc/internal/slots"
)

func Print(s slots.Slots, trace *rules.ApplyResult) {
	fmt.Print("Explain mode\n\n")

	printMeta(s)
	printIntent(s)
	printTopic(s)
	printStage(s)
	printEntities(s)
	printModifiers(s)

	if trace != nil {
		printRuleTrace(trace)
	}
}

func printMeta(s slots.Slots) {
	fmt.Println("Meta:")
	fmt.Printf("  Language: %s\n\n", valueOrNone(s.Language))
}

func printIntent(s slots.Slots) {
	fmt.Println("Intent:")
	fmt.Printf("  Detected intent: %s\n\n", valueOrNone(s.Intent))
}

func printTopic(s slots.Slots) {
	fmt.Println("Topic:")
	if s.Topic == "" {
		fmt.Print("  <none>\n\n")
		return
	}
	fmt.Printf("  \"%s\"\n\n", s.Topic)
}

func printStage(s slots.Slots) {
	fmt.Println("Stage:")
	fmt.Printf("  %s\n\n", valueOrNone(s.Stage))
}

func printEntities(s slots.Slots) {
	fmt.Println("Entities:")

	if len(s.Entities) == 0 {
		fmt.Print("  <none>\n\n")
		return
	}

	for _, e := range s.Entities {
		fmt.Printf("  - %s: %s\n", e.Role, e.Text)
	}
	fmt.Println()
}

func printModifiers(s slots.Slots) {
	fmt.Println("Modifiers:")

	fmt.Printf("  Audience: %s\n", valueOrNone(s.Audience))
	fmt.Printf("  Depth:    %s\n", valueOrNone(s.Depth))
	fmt.Printf("  Style:    %s\n", valueOrNone(s.Style))
	fmt.Printf("  Format:   %s\n\n", valueOrNone(s.Format))
}

func printRuleTrace(trace *rules.ApplyResult) {
	fmt.Println("Rules:")
	fmt.Printf("  Applied: %s\n", formatRuleList(trace.Applied))
	fmt.Printf("  Skipped: %s\n", formatRuleList(trace.Skipped))
}

func formatRuleList(ids []string) string {
	if len(ids) == 0 {
		return "<none>"
	}
	return strings.Join(ids, ", ")
}

func valueOrNone(v string) string {
	if v == "" {
		return "<none>"
	}
	return v
}
