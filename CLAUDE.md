# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run a single test by name
go test -run TestConcat ./hiter

# Run tests for a specific package
go test ./hiter/stringsiter/...

# Run example tests
go test -run Example ./...
```

### Code Quality
```bash
# Format all code
go fmt ./...

# Run static analysis
go vet ./...

# Ensure module dependencies are correct
go mod tidy
```

## Architecture Overview

### Iterator Foundation
This library builds on Go 1.23's iterator types (`iter.Seq[T]` for single values and `iter.Seq2[K,V]` for key-value pairs). Every function is designed to either:
- Generate iterators from various sources (slices, channels, ranges)
- Transform existing iterators (map, filter, flatten)
- Consume iterators (reduce, collect, sum)

### Package Organization Pattern
The codebase follows a deliberate mirroring strategy:
- `hiter/` - Core iterator functionality, avoiding duplication with standard library
- `hiter/*iter/` - Each sub-package corresponds to a standard library package:
  - `stringsiter` → `strings` package iterators
  - `reflectiter` → `reflect` package iterators
  - `encodingiter` → Flattened encoding/* packages (json, xml, csv)
  - `databaseiter` → `database/sql` iterators

### Testing Architecture
Tests use a custom `testcase` framework (`hiter/internal/testcase`) that provides:
- `testcase.One[T]` for testing `iter.Seq[T]` iterators
- `testcase.Two[K,V]` for testing `iter.Seq2[K,V]` iterators
- Automatic validation of both complete iteration and early break scenarios
- Example tests demonstrate real-world usage patterns

### Error Handling Pattern
The library uses `errbox` package for error propagation in iterators since iterators can't return errors directly. This allows errors to be captured during iteration and checked afterward.

### Key Design Principles
1. **No Reinvention**: Never implement what exists in standard library (`slices`, `maps`)
2. **Dual API Pattern**: Most functions have both `Seq` and `Seq2` variants (e.g., `Map` and `Map2`)
3. **Lazy Evaluation**: All iterators are lazy; computation happens during iteration
4. **Composition Focus**: Small, composable functions that can be chained together

### Function Naming Conventions
- Simple names for `Seq[T]` functions: `Map`, `Filter`, `Reduce`
- Suffix `2` for `Seq2[K,V]` variants: `Map2`, `Filter2`, `Reduce2`
- Past tense for collection functions: `Collect`, `Collected`
- `*All` suffix for complete iteration variants: `SliceAll`, `FindAll`

## Important Notes from README

- The library avoids implementing iterator sources for types that are already iterators (like `sync.Map.Range`)
- Functions are deprecated when standard library equivalents become available, but remain functional
- The deprecated `x/exp/xiter` package was withdrawn; `hiter` now contains those implementations

## Working with Code

**MUST** use Serena tools for code exploration and modification whenever possible:
- **MUST** use `mcp__serena__get_symbols_overview` and `mcp__serena__find_symbol` to explore code structure
- **MUST** use `mcp__serena__search_for_pattern` to search for patterns across files
- **MUST** use `mcp__serena__replace_symbol_body` or `mcp__serena__insert_*_symbol` for code modifications
- **MUST NOT** read entire files when only specific symbols are needed
- **MUST NOT** use basic text editing when semantic symbol-based editing is available