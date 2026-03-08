# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-03-08

### Changed
- Unified language files: translations and language metadata now live in a single file per language (`languages/{en,de}.yaml`)
- Removed `translations/` directory — all content merged into `languages/`
- German language name changed from "German" to "Deutsch"

## [1.2.0] - 2026-02-25

### Added
- Tiered enrichment system: prompts automatically gain semantic depth (role, context, scope, constraints) based on input richness
- Three enrichment tiers: minimal (base rules only), standard (moderate enrichments), rich (full enrichments across all sections)
- Cross-field interaction rules that adapt output based on slot combinations (e.g., beginner × explain, advanced × deep)
- Partial credit scoring: list sections score per-item instead of all-or-nothing
- Full bidirectional translation symmetry validation for EN and DE

### Changed
- Richer prompt output for standard and rich tier inputs
- Scoring reflects incremental completeness more accurately

### Fixed
- Non-deterministic phrase ordering that could cause subtle output variations across runs

## [1.1.0] - 2026-02-25

### Added
- Cascading topic extraction with three-tier fallback (phrase matching, boundary-aware, simple split)
- Multi-word phrase dictionary (`data/phrases.yaml`) with per-language EN/DE entries
- Three new intents: `debug`, `refactor`, `summarize` with dedicated scope rules
- Per-section score breakdown with bar charts (`--score`)
- Stdin piping: `echo "input" | promptc compile` and `promptc compile -`
- Full history management: `--search`, `--delete`, `--clear`, `--export`, `--limit` flags
- REPL `:search` command for history search
- Auto-pruning: history capped at 500 entries / 90 days
- Locale-aware time formatting via `goodsign/monday` (driven by translation YAML)
- Expanded acronyms dictionary (~200 entries across 12 categories)

### Changed
- Rule count increased from 25 to 28
- Score display shows per-section breakdown instead of just total
- Time formats in history/REPL now respect language setting

## [1.0.1] - 2026-02-22

### Fixed
- CI workflow compatibility with golangci-lint v2 config schema
- Updated GitHub Actions versions (checkout v5, setup-go v6, golangci-lint-action v9)

## [1.0.0] - 2026-02-22

### Added
- Deterministic prompt compilation from natural language input
- Five intent types: explain, howto, generate, analyze, decide
- Rule-based NLP pipeline with 25 built-in rules
- Interactive REPL mode with colored output and commands
- JSON and YAML output formats
- Prompt completeness scoring (0-100)
- Bilingual support (English and German)
- Language detection via fastText with keyword fallback
- Clipboard integration
- Prompt history with recall
- Debug mode with rule tracing (--explain)
- Shell completions (bash, zsh, fish, powershell)
