# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
