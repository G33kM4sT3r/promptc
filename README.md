# promptc

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/G33kM4sT3r/promptc/actions/workflows/ci.yml/badge.svg)](https://github.com/G33kM4sT3r/promptc/actions/workflows/ci.yml)
[![Offline](https://img.shields.io/badge/Offline-100%25-brightgreen)]()
[![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-lightgrey)]()
[![Buy Me A Coffee](https://img.shields.io/badge/Buy_Me_A_Coffee-FFDD00?logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/martin.willig)

A deterministic, offline prompt compiler. Transforms informal natural language into structured, high-quality AI prompts using rule-based NLP — no LLMs, no API calls, no cloud dependencies.

## Install

**From source** (requires Go 1.26+ and `make`):

```bash
git clone https://github.com/G33kM4sT3r/promptc.git
cd promptc
make download-model  # optional: fastText language detection
make install
```

Or download a [source release](https://github.com/G33kM4sT3r/promptc/releases), extract, and run `make install`.

CGO must be enabled (the default) to use the optional fastText language detection model. Without it, keyword-based detection is used automatically.

## Usage

```bash
promptc compile "explain closures for beginners"
```

Output:

```
knowledgeable instructor

Objective:
Explain closures.

Context:
The explanation should be beginner-friendly.
Start with the big picture before diving into details.

Scope:
- Explain the main concept
- Cover the most important aspects
- Use an analogy or real-world comparison to build intuition
- Start with the big picture before details

Constraints:
- Use simple and clear language
- Avoid unnecessary jargon

Output:
- Use clear section headings

Quality criteria:
- Clear and structured
- Accurate
```

### Intents

```bash
promptc compile "explain closures for beginners"           # explain
promptc compile "how do I start a project with PHP"         # howto
promptc compile "generate a REST API with Python"           # generate
promptc compile "analyze authentication code with Java"     # analyze
promptc compile "should I use React or Vue"                 # decide
promptc compile "debug why my API returns 500"              # debug
promptc compile "refactor this function for readability"    # refactor
promptc compile "summarize microservices architecture"      # summarize
```

### Output Formats

Choose between text, JSON, or YAML output with `--output` (or `-o`):

```bash
promptc compile --output json "explain closures for beginners"
promptc compile --output yaml "generate a REST API with Python"
promptc compile -o json "analyze authentication code"
```

### Quality Score

Use `--score` to see a per-section completeness breakdown (0–100) with bar charts:

```bash
promptc compile --score "explain closures for beginners"
```

The breakdown is printed to stderr so it doesn't interfere with piped output. Also available as `--score --output json` or `--score --output yaml`.

### Clipboard

Copy the compiled prompt directly to your clipboard with `--copy`:

```bash
promptc compile --copy "generate a REST API with Python"
```

### Debug Mode

Use `--explain` (or `-e`) to see how the input was interpreted, including rule tracing:

```bash
promptc compile --explain "how do I start a project with PHP"
```

### German Support

German input is automatically detected and all output is rendered in German — objectives, scope, constraints, and labels:

```bash
promptc compile "erkläre dependency injection detailliert"
```

Output:

```
Experte, der Erklärungen an das Publikum anpasst

Ziel:
Erkläre dependency injection.

Kontext:
Berücksichtige, welches Grundlagenwissen das Publikum benötigt.
Baue Verständnis schrittweise von einfach zu komplex auf.

Umfang:
- Die Kernkonzepte im Detail erklären
- Praktische Beispiele liefern
- Wichtige Nuancen besprechen
- Konkrete Beispiele liefern, die jedes Konzept veranschaulichen
- Häufige Missverständnisse ansprechen
- Grenzfälle besprechen und wann das Konzept an seine Grenzen stößt

Einschränkungen:
- Gib detaillierte Erklärungen
- Fachbegriffe definieren, bevor sie verwendet werden

Ausgabeformat:
- Verwende klare Abschnittsüberschriften
- Jeden Abschnitt mit einer einzeiligen Zusammenfassung beginnen

Qualitätskriterien:
- Klar und strukturiert
- Korrekt
- Zugänglich, ohne an Korrektheit einzubüßen
```

### Interactive Mode

Start an interactive REPL for iterative prompt compilation:

```bash
promptc repl
```

The REPL provides a persistent session with colored output, a status bar showing the current language and explain mode, and built-in commands:

| Command              | Description                          |
|----------------------|--------------------------------------|
| `:help`, `:h`        | Show available commands              |
| `:explain`           | Toggle explain mode (rule tracing)   |
| `:lang <code>`       | Set language (e.g., `:lang de`)      |
| `:lang`              | Show current language                |
| `:output <fmt>`      | Switch output format (text/json/yaml)|
| `:history`           | List previous compilations           |
| `:recall <N>`        | Recall a history entry by index      |
| `:search <term>`     | Search history by input text         |
| `:copy`              | Copy last output to clipboard        |
| `:quit`, `:q`, `:exit` | Exit the REPL                     |

The REPL automatically scores every compilation and saves it to history. Use `Ctrl+C` or `Ctrl+D` to quit.

### Language Override

Use `--lang` (or `-l`) to force a specific language regardless of detection:

```bash
promptc compile --lang de "explain closures for beginners"   # render in German
promptc compile --lang en "erkläre closures"                 # render in English
```

### Stdin Piping

Pipe input from other commands instead of passing it as an argument:

```bash
echo "explain closures" | promptc compile
cat prompt.txt | promptc compile --output json
promptc compile - < prompt.txt
```

### Prompt History

Compilations from both the CLI and REPL are saved automatically (capped at 500 entries / 90 days):

```bash
promptc history                       # list all entries
promptc history --limit 5             # show last 5
promptc history 0                     # recall entry by index
promptc history 0 --output json       # recall as JSON
promptc history --search "explain"    # search by input text
promptc history --delete 2 --yes      # delete entry at index
promptc history --clear --yes         # clear all history
promptc history --export              # export as JSON
```

### Version

```bash
promptc version
```

## Language Detection

promptc detects language automatically. For higher accuracy, download the fastText model:

```bash
make download-model
```

This downloads `data/lid.176.ftz` (~900 KB). When the model is absent, promptc falls back to a keyword-scoring heuristic that reliably distinguishes English and German.

## How It Works

promptc compiles prompts through a pipeline of deterministic stages:

```
Input → Normalize → Tokenize (phrase → boundary → split) → Detect Language → Extract Slots → Clean Topic → Apply Rules → Render
```

**Config** loads comprehensive keyword lexicons from YAML data files at startup. All extraction is data-driven — extend keywords without changing code.

**Language Detection** uses fastText (`data/lid.176.ftz`) when available, falling back to keyword-score heuristics. The `--lang` (`-l`) flag overrides detection entirely.

**Extraction** identifies intent (`explain`, `howto`, `generate`, `analyze`, `decide`, `debug`, `refactor`, `summarize`), topic, entities (e.g., "with PHP" as implementation medium), stage (`getting-started`, `implementation`, `optimization`), and modifiers (audience, depth, style, format) using config-driven lookups with phrase and keyword matching. A **tier** (`minimal`, `standard`, or `rich`) is then computed from slot richness (depth, audience, entity count, stage).

**Topic Cleanup** strips leading articles ("a", "the", "an", "die", "der", "das") and normalizes known acronyms to their canonical casing (`data/acronyms.yaml`).

**Rules** (34 built-in) transform extracted slots into a structured `PromptSpec` — encoding prompt engineering best practices as composable, order-aware functions. Each rule accepts a `Translator` to produce language-appropriate text. Tier-based enrichment rules add semantic depth (role, context, scope, constraints) based on the input's richness, and cross-field interaction rules handle slot combinations like audience×intent or stage×depth.

**Rendering** formats the spec into text, JSON, or YAML. The `TranslatedRenderer` looks up section labels (`Objective:`, `Ziel:`, etc.) from the active translation file. JSON and YAML renderers serialize the spec directly for machine consumption.

## Extending

- **Keywords** — add to `data/*.yaml` and `languages/*.yaml` (no code changes needed)
- **Languages** — add a new `translations/<lang>.yaml` file with translated label and rule strings
- **Rules** — add to `internal/rules/builtin/` and register in the pipeline
- **Acronyms** — add to `data/acronyms.yaml` for canonical casing in topics
- **Enrichments** — edit `data/enrichments.yaml` to add tier-based content per intent/section; actual text lives in `translations/*.yaml` under `enrichment.*` and `cross.*` namespaces

## Design Principles

- **Deterministic**: same input always produces the same output
- **Offline**: no network calls, no API keys, no cloud services
- **Transparent**: `compile --explain` shows every extraction decision and which rules fired
- **Conservative**: prefers empty slots over incorrect guesses
- **Data-driven**: extend via YAML config, not code changes
- **i18n by files**: add a language by adding a YAML translation file — no rule or renderer code changes

## License

MIT License – see [LICENSE](LICENSE) for details.
