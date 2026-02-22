package pipeline

import (
	"errors"
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
	tokens := tokenize.Tokenize(norm)

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
	return s, nil
}

// ApplyRules applies the rule engine to pre-extracted slots.
func ApplyRules(s slots.Slots, t *i18n.Translator) prompt.PromptSpec {
	engine := newEngine()
	return engine.Apply(s, t)
}

// ApplyRulesWithTrace applies the rule engine with tracing to pre-extracted slots.
func ApplyRulesWithTrace(s slots.Slots, t *i18n.Translator) rules.ApplyResult {
	engine := newEngine()
	return engine.ApplyWithTrace(s, t)
}

// RunWithRules processes input through extraction and then applies rules to build a PromptSpec.
// If det is nil, a keyword-only detector is created automatically.
func RunWithRules(input string, cfg *config.Config, t *i18n.Translator, det *language.Detector) (slots.Slots, prompt.PromptSpec, error) {
	s, err := Run(input, cfg, det)
	if err != nil {
		return s, prompt.PromptSpec{}, err
	}

	spec := ApplyRules(s, t)
	return s, spec, nil
}

// RunWithTrace is like RunWithRules but also returns rule tracing information.
// If det is nil, a keyword-only detector is created automatically.
func RunWithTrace(input string, cfg *config.Config, t *i18n.Translator, det *language.Detector) (slots.Slots, rules.ApplyResult, error) {
	s, err := Run(input, cfg, det)
	if err != nil {
		return s, rules.ApplyResult{}, err
	}

	result := ApplyRulesWithTrace(s, t)
	return s, result, nil
}

func newEngine() *rules.Engine {
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
	})
}
