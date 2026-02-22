package config

// Config holds all loaded configuration data for the prompt compiler.
type Config struct {
	Intents      IntentsConfig             `yaml:"intents"`
	Modifiers    ModifiersConfig           `yaml:"modifiers"`
	Stages       map[string]StageEntry     `yaml:"stages"`
	Entities     map[string]EntityEntry    `yaml:"entities"`
	Languages    map[string]LanguageConfig `yaml:"languages"`
	Acronyms     AcronymsConfig            `yaml:"-"`
	Contractions ContractionsConfig        `yaml:"-"`
}

// ContractionsConfig holds contraction expansion mappings.
type ContractionsConfig struct {
	Contractions map[string]string `yaml:"contractions"`
}

// AcronymsConfig holds the list of known acronyms for topic casing.
type AcronymsConfig struct {
	Acronyms []string `yaml:"acronyms"`
}

// IntentsConfig maps intent names to their detection patterns.
type IntentsConfig struct {
	Intents map[string]IntentEntry `yaml:"intents"`
}

// IntentEntry defines keywords and phrases for a single intent.
type IntentEntry struct {
	Keywords map[string][]string `yaml:"keywords"` // lang → keyword list
	Phrases  map[string][]string `yaml:"phrases"`  // lang → phrase list
}

// ModifiersConfig holds all modifier categories.
type ModifiersConfig struct {
	Audience map[string]ModifierEntry `yaml:"audience"`
	Depth    map[string]ModifierEntry `yaml:"depth"`
	Style    map[string]ModifierEntry `yaml:"style"`
	Format   map[string]ModifierEntry `yaml:"format"`
}

// ModifierEntry maps a modifier value to its keywords per language.
type ModifierEntry struct {
	Keywords map[string][]string `yaml:"keywords"` // lang → keyword list
}

// StageEntry defines keywords for a lifecycle stage.
type StageEntry struct {
	Keywords map[string][]string `yaml:"keywords"` // lang → keyword list
}

// EntityEntry defines an entity role's detection patterns.
type EntityEntry struct {
	Prepositions map[string][]string `yaml:"prepositions"` // lang → preposition list
	MultiWord    bool                `yaml:"multi_word"`
	StopWords    map[string][]string `yaml:"stop_words"` // lang → stop word list
}

// LanguageConfig defines per-language settings.
type LanguageConfig struct {
	Name       string   `yaml:"name"`
	Code       string   `yaml:"code"`
	StopWords  []string `yaml:"stop_words"`
	TopicVerbs []string `yaml:"topic_verbs"`
}
