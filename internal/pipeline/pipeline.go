package pipeline

import (
	"errors"
	"sort"
	"strings"

	"promptc/internal/config"
	"promptc/internal/extract"
	"promptc/internal/i18n"
	"promptc/internal/language"
	"promptc/internal/normalize"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/rules/builtin"
	"promptc/internal/slots"
	"promptc/internal/tokenize"
)

// Run processes input through the extraction pipeline and returns populated Slots.
// If det is nil, a keyword-only detector is created automatically.
func Run(input string, cfg *config.Config, det *language.Detector) (slots.Slots, error) {
	if strings.TrimSpace(input) == "" {
		return slots.Slots{}, errors.New("empty input")
	}

	if det == nil {
		var err error
		det, err = language.NewDetector("", cfg)
		if err != nil {
			return slots.Slots{}, err
		}
	}

	norm := normalize.NormalizeAll(input, cfg.Contractions.Contractions)
	tokens := tokenize.TokenizeWithPhrases(norm, allPhrases(cfg.Phrases.Phrases))

	ext := extract.New(cfg)

	s := slots.Slots{
		Language: det.Detect(norm),
		Intent:   ext.DetectIntent(tokens),
		Topic:    ext.ExtractTopic(tokens),
		Stage:    ext.DetectStage(tokens),
		Entities: ext.ExtractEntities(tokens),
		Audience: ext.DetectAudience(tokens),
		Depth:    ext.DetectDepth(tokens),
		Style:    ext.DetectStyle(tokens),
		Format:   ext.DetectFormat(tokens),
	}
	s.Topic = ext.CleanTopic(s.Topic)
	s.Tier = extract.CalculateTier(s)
	return s, nil
}

// ApplyRules applies the rule engine to pre-extracted slots.
func ApplyRules(s slots.Slots, t *i18n.Translator, enrichments config.EnrichmentsConfig) prompt.PromptSpec {
	engine := newEngine(enrichments)
	return engine.Apply(s, t)
}

// ApplyRulesWithTrace applies the rule engine with tracing to pre-extracted slots.
func ApplyRulesWithTrace(s slots.Slots, t *i18n.Translator, enrichments config.EnrichmentsConfig) rules.ApplyResult {
	engine := newEngine(enrichments)
	return engine.ApplyWithTrace(s, t)
}

// RunWithRules processes input through extraction and then applies rules to build a PromptSpec.
// If det is nil, a keyword-only detector is created automatically.
func RunWithRules(input string, cfg *config.Config, t *i18n.Translator, det *language.Detector) (slots.Slots, prompt.PromptSpec, error) {
	s, err := Run(input, cfg, det)
	if err != nil {
		return s, prompt.PromptSpec{}, err
	}

	spec := ApplyRules(s, t, cfg.Enrichments)
	return s, spec, nil
}

// RunWithTrace is like RunWithRules but also returns rule tracing information.
// If det is nil, a keyword-only detector is created automatically.
func RunWithTrace(input string, cfg *config.Config, t *i18n.Translator, det *language.Detector) (slots.Slots, rules.ApplyResult, error) {
	s, err := Run(input, cfg, det)
	if err != nil {
		return s, rules.ApplyResult{}, err
	}

	result := ApplyRulesWithTrace(s, t, cfg.Enrichments)
	return s, result, nil
}

// allPhrases merges per-language phrase lists into a single deduplicated list.
// Keys are sorted to ensure deterministic output regardless of map iteration order.
func allPhrases(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	seen := make(map[string]bool)
	var result []string
	for _, k := range keys {
		for _, p := range m[k] {
			if !seen[p] {
				seen[p] = true
				result = append(result, p)
			}
		}
	}
	return result
}

func newEngine(enrichments config.EnrichmentsConfig) *rules.Engine {
	return rules.NewEngine([]rules.Rule{
		builtin.ObjectiveRule(),
		builtin.RoleFromEntitiesRule(),

		// Context
		builtin.ContextFromStageRule(),
		builtin.ContextFromAudienceRule(),
		builtin.ContextFromTargetObjectRule(),

		// Scope (specific rules before fallback)
		builtin.ScopeHowToGettingStartedRule(),
		builtin.ScopeHowToImplementationRule(),
		builtin.ScopeHowToOptimizationRule(),
		builtin.ScopeExplainRule(),
		builtin.ScopeGenerateRule(),
		builtin.ScopeAnalyzeRule(),
		builtin.ScopeDecideRule(),
		builtin.ScopeDebugRule(),
		builtin.ScopeRefactorRule(),
		builtin.ScopeSummarizeRule(),
		builtin.ScopeFromExampleObjectRule(),
		builtin.ScopeFallbackRule(),

		// Constraints
		builtin.ConstraintsFromDepthRule(),
		builtin.ConstraintsFromAudienceRule(),
		builtin.ConstraintsFromStyleRule(),
		builtin.ConstraintsFromEntitiesRule(),
		builtin.ConstraintsFromConstraintObjectRule(),

		// Output
		builtin.OutputFromFormatRule(),
		builtin.OutputFromIntentRule(),
		builtin.OutputFromDepthRule(),
		builtin.OutputFallbackRule(),

		// Quality
		builtin.QualityBaseRule(),
		builtin.QualityFromIntentRule(),

		// Tier enrichment
		builtin.EnrichFromTierRule(enrichments),

		// Cross-field interactions
		builtin.CrossAudienceIntentRule(),
		builtin.CrossEntityIntentRule(),
		builtin.CrossStageDepthRule(),
		builtin.CrossAudienceDepthRule(),
		builtin.CrossStyleAudienceRule(),
	})
}
