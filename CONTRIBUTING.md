# Contributing to promptc

## Development Setup

1. Clone the repository
2. Ensure Go 1.26+ is installed
3. Run `make build` to compile
4. Run `make test` to verify everything works

Optionally download the fastText language model:

```bash
make download-model
```

## Running Tests

```bash
make test          # Run all tests with race detector
make cover         # Generate coverage report
make cover-html    # Open coverage in browser
```

## Code Quality

Before submitting changes:

```bash
make check         # Runs fmt-check + vet + lint + test
```

This requires [golangci-lint](https://golangci-lint.run/usage/install/) to be installed.

## Pull Requests

- Keep changes focused and minimal
- Ensure all tests pass and coverage stays above 80%
- Follow existing code conventions (run `make fmt` before committing)
- Update translations in both `translations/en.yaml` and `translations/de.yaml` if adding user-facing strings
- Add tests for new functionality

## Project Structure

See the architecture section in [README.md](README.md) or [CLAUDE.md](CLAUDE.md) for detailed package descriptions.
