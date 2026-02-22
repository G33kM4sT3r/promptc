# CLAUDE.md

promptc is a deterministic, offline CLI prompt compiler (Go, CGO required for `go-fasttext`). Transforms natural language → structured AI prompts via rule-based NLP. No LLMs, no API calls, no randomness.

## Commands

```bash
make build                    # Build binary (CGO_ENABLED=1)
make test                     # Tests with -race
make check                    # fmt-check + vet + lint + test
make cover                    # Coverage report (excludes cmd/)
make download-model           # Download fastText lid.176.ftz → data/
go run ./cmd/promptc compile "input"              # Compile prompt
go run ./cmd/promptc compile --explain "input"    # Rule tracing
go run ./cmd/promptc compile --lang de "input"    # Language override
go run ./cmd/promptc compile --output json "input" # JSON/YAML output
go run ./cmd/promptc compile --score "input"      # Completeness score
go run ./cmd/promptc compile --copy "input"       # Copy to clipboard
go run ./cmd/promptc repl                         # Interactive REPL
go run ./cmd/promptc history [N]                  # List/recall history
go run ./cmd/promptc version                      # Version info
go run ./cmd/promptc completion bash              # Shell completions
```

## Architecture

Linear pipeline: `Input → Normalize → Tokenize → Detect Language → Extract Slots → Clean Topic → Apply Rules (+ Translator) → Render → Output`

| Package | Purpose |
|---|---|
| `config/` | Loads YAML from `data/` + `languages/`. `FindBaseDir()` checks `PROMPTC_DATA` env, exe dir, cwd. |
| `language/` | fastText ML detection → keyword fallback → ISO 639-1 code. `--lang` overrides. |
| `pipeline/` | Orchestrates stages: `Run()` → Slots, `RunWithRules()` → PromptSpec, `RunWithTrace()` → explain info. |
| `extract/` | Config-driven slot extraction via reverse lookup maps. Intent (5 types), topic, stage, entities (4 roles), modifiers (audience/depth/style/format). |
| `i18n/` | `Translator` — keyed lookup from `translations/<lang>.yaml`, English fallback, key-string last resort. `Get(key)`, `Getf(key, args...)`. |
| `rules/` | 25 order-dependent, append-only rules in `builtin/`. `When` guard + `Apply` mutator. Deduplicates after all rules fire. `ApplyWithTrace()` for explain mode. |
| `repl/` | Bubbletea TUI. Commands: `:help :explain :lang :output :history :recall :copy :quit`. Auto-saves to history. |
| `score/` | Completeness 0–100. Weights: objective(25), context(15), scope(15), output(15), role(10), constraints(10), quality(10). |
| `history/` | `Store` persists to `history.json`. API: `Add()`, `List()`, `Get(index)`. |
| `clipboard/` | Wraps `atotto/clipboard`. CLI: `--copy`. REPL: `:copy`. |
| `ui/` | Color palette, gradients, `NO_COLOR`/TTY detection, bubbletea spinner. |
| `render/` | `Renderer` interface: `TranslatedRenderer` (styled text), `JSONRenderer`, `YAMLRenderer`. |

### Key Types

- `slots.Slots` — linguistic signals (intent, topic, stage, entities, modifiers, language)
- `prompt.PromptSpec` — structured sections (role, objective, context, scope, constraints, output, quality)
- `rules.Rule` — `{ID, When func(Slots) bool, Apply func(*PromptSpec, Slots, *Translator)}`
- `config.Config` — all loaded YAML data
- `score.ScoreResult` — `{Total int, Breakdown map[string]int}`

### Rule Ordering

Rules are order-dependent — `rule_ordering_test.go` encodes constraints. When adding rules: add to `pipeline.newEngine()` AND the test's `canonicalRuleOrder()`.

## Design Principles

- **Deterministic**: same input + config = identical output
- **Conservative inference**: prefer empty slots over guessing
- **Data-driven**: extend via `data/*.yaml`, `languages/*.yaml`, `translations/*.yaml` — no code changes
- **Bilingual**: extraction handles EN + DE keywords; rendering fully translated via Translator

## File Conventions

- CLI: `cmd/promptc/main.go` (minimal), commands in `cmd/promptc/{root,compile,repl,history,version}.go`
- Rules: `internal/rules/builtin/` named by category (`scope_explain.go`, `constraints.go`)
- Data: `data/*.yaml`, `languages/*.yaml`, `translations/{en,de}.yaml`
- Tests: golden files in `internal/pipeline/testdata/*.golden`, `doc.go` per package
- Runtime: `history.json` (gitignored), `data/lid.176.ftz` (optional, gitignored)
- Worktrees: `.worktrees/` (gitignored)

## Known Gotchas

- **macOS linker warning**: `ld: warning: ignoring duplicate libraries: '-lc++'` from go-fasttext CGO. Suppressed in Makefile via `CGO_LDFLAGS`. Harmless.
- **Extraction case**: Lookup maps use lowercase keys. Always `strings.ToLower()` before lookups, preserve original case in output.
- **Translation keys**: Both `en.yaml` and `de.yaml` must have symmetric keys. No orphaned keys.
- **golangci-lint v2**: Requires `version: "2"`. `gofmt` under `formatters:` not `linters:`. `gosimple` merged into `staticcheck`. `errcheck.check-type-assertions` ignores path-based excludes.
- **Coverage**: `cmd/promptc/` excluded — integration tests use `exec.Command` on compiled binary. Use `go list ./... | grep -v promptc/cmd/` for coverage. CLI tests run separately.
- **CLI integration tests**: Require `PROMPTC_DATA` env var set to project root for binary to find `data/` and `languages/`.
- **UI color tests**: Non-TTY = no ANSI. `NO_COLOR` uses `os.LookupEnv` (existence check). Set `colorsOff = false` directly in tests.
