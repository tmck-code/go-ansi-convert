# go-ansi-convert AI Agent Instructions

## Project Overview
A Go CLI tool for manipulating ANSI art (colored ASCII art). Core operations: **flip** (horizontal/vertical), **sanitize** (normalize ANSI codes), **justify** (align lines), and **optimize** (merge redundant color codes).

## Architecture & Core Concepts

### Token-Based Processing Model
All ANSI operations work on **tokenized representations**, not raw strings:
- Raw ANSI string → `TokeniseANSIString()` → `[][]ANSILineToken` → operations → `BuildANSIString()` → output
- Each `ANSILineToken` has: `FG` (foreground), `BG` (background), `T` (text content)
- ANSI codes are separated from text for independent manipulation
- See [src/convert/convert.go](src/convert/convert.go) lines 112-320 for tokenization logic

#### Tokenization Algorithm Details
The tokenizer is a **character-by-character state machine**:
- Detects ANSI escape sequences starting with `\x1b[` (ESC + `[`)
- Accumulates characters in `colour` buffer until `m` terminator is found
- Handles multiple ANSI formats:
  - **256-color codes**: `\x1b[38;5;Nm` (FG) / `\x1b[48;5;Nm` (BG) - detected by `;5;` pattern
  - **True color codes**: `\x1b[38;2;R;G;Bm` - detected by `;2;` pattern
  - **8/16-color codes**: `\x1b[30-37m`, `\x1b[40-47m`, `\x1b[90-97m`, `\x1b[100-107m`
  - **Combined codes**: `\x1b[0;31;40m` - splits into separate FG/BG tokens
  - **Reset codes**: `\x1b[0m` - clears both FG and BG state
- **Background reset insertion**: Automatically adds `\x1b[49m` to previous token when BG changes
- Text between ANSI codes accumulates in `text` buffer, flushed when new ANSI code starts
- State persists across each line but resets between lines

### Unicode Width Awareness
The project handles **double-width characters** (CJK, emojis, box-drawing chars):
- Use `UnicodeStringLength()` instead of `len()` for display width calculations
- Uses `github.com/mattn/go-runewidth` to calculate actual terminal width
- Critical for flipping/justifying: characters like `▄▀█` may have width > 1
- Example in [src/convert/convert.go](src/convert/convert.go) lines 9-38

### Character Mirroring Maps
Horizontal/vertical flips include **character substitution** via lookup tables:
- `HorizontalMirrorMap`: `'<'` ↔ `'>'`, `'d'` ↔ `'b'`, `'/'` ↔ `'\'`, etc.
- `VerticalMirrorMap`: `'▀'` ↔ `'▄'`, etc.
- Maps are bidirectional (both directions defined)
- See [src/convert/mirror.go](src/convert/mirror.go) for complete mappings

### Color State Management
ANSI color tracking is stateful and complex:
- Background reset (`\x1b[49m`) is automatically inserted between tokens when needed
- Reset codes (`\x1b[0m`) clear both FG and BG state
- 256-color codes: `\x1b[38;5;Nm` (FG) and `\x1b[48;5;Nm` (BG)
- Tokenizer maintains `fg`, `bg`, and `isReset` state across line parsing

## Testing Patterns

### Test Structure
- Tests live in `test/*` subdirectories (not `src/convert/*_test.go`)
- Use table-driven tests with `expected [][]convert.ANSILineToken` format
- Run with `DEBUG=true go test -v ./test/...` to see token-level comparisons

### Custom Test Helpers
Located in [test/helper.go](test/helper.go):
- `Assert(expected, result, t)`: compares `%#v` representations
- `PrintANSITestResults()`: renders input/expected/result with visual borders
- `Debug()`: checks `DEBUG=true` env var for verbose output
- Visual markers: `✓` (success), `❌` (fail), colored titles for sections

### Integration Test Data
Complex ANSI sprites in `test/data/`:
- Pikachu, egg sprites with full 256-color ANSI codes
- Used for flip integration tests (see [test/flip/flip_test.go](test/flip/flip_test.go) lines 88-145)

## Development Workflows

### Running Tests
```bash
# Run all tests
go test -v ./test/...

# With debug output (shows token diffs)
DEBUG=true go test -v ./test/...

# Docker-based testing (matches CI)
docker build -t go-ansi-convert-test -f test/Dockerfile .
docker run --rm go-ansi-convert-test go test -v ./test/...
```

### Building & Running
```bash
# Build binary
go build -o ansi-flip main.go

# Example usage
./ansi-flip --flip h --input sprite.ans
./ansi-flip --sanitise --justify < input.ans > output.ans
```

### Writing any tmp/debug files

Write any tmp/debug files to this repo, not /tmp/

### CI Pipeline
GitHub Actions runs tests in Docker (see [.github/workflows/test.yml](.github/workflows/test.yml))
- Builds `test/Dockerfile`
- Runs full test suite on push

## Common Operations & Examples

### Horizontal Flip
```go
tokens := convert.TokeniseANSIString(input)
flipped := convert.FlipHorizontal(tokens)
output := convert.BuildANSIString(flipped, 0)
```
- Reverses token order on each line (right-to-left)
- Reverses characters within each token's text
- Applies character mirroring: `d` → `b`, `<` → `>`, `/` → `\`
- Pads left side with spaces for vertical alignment (based on widest line)

### Vertical Flip
```go
tokens := convert.TokeniseANSIString(input)
flipped := convert.FlipVertical(tokens)
output := convert.BuildANSIString(flipped, 0)
```
- Reverses line order (bottom-to-top)
- Applies character mirroring: `▀` → `▄`
- Preserves horizontal ordering within lines

### Sanitize with Justification
```go
sanitised := convert.SanitiseUnicodeString(input, true) // justify=true
```
- Ensures all lines end with `\x1b[0m` reset code
- Pads lines to equal width (based on `UnicodeStringLength`)
- Preserves ANSI colors while normalizing structure

### Optimize ANSI Tokens
```go
tokens := convert.TokeniseANSIString(input)
optimised := convert.OptimiseANSITokens(tokens)
output := convert.BuildANSIString(optimised, 0)
```
- **Merges adjacent tokens** with identical FG/BG colors
- **Eliminates redundant codes**: Only outputs FG when it changes, only BG when it changes
- **Removes empty reset tokens**: `{FG: "\x1b[0m", T: ""}` with no text
- Example: `{FG: "\x1b[38;5;129m", BG: "", T: "A"}` + `{FG: "\x1b[38;5;129m", BG: "", T: "B"}` → `{FG: "\x1b[38;5;129m", BG: "", T: "AB"}`
- Reduces output size and terminal rendering overhead

### Full Processing Pipeline (from main.go)
```go
// 1. Read input
input := readStdin() // or readFile(path)

// 2. Tokenize
tokens := convert.TokeniseANSIString(input)

// 3. Apply operations
tokens = convert.FlipHorizontal(tokens)  // optional
tokens = convert.FlipVertical(tokens)    // optional
optimised := convert.OptimiseANSITokens(tokens)  // optional

// 4. Rebuild output
output := convert.BuildANSIString(optimised, 0)
```

## Code Conventions

### Naming: British English
- Use `Sanitise`, `Optimise`, `Tokenise` (not `Sanitize`, etc.)
- Consistent across function names and documentation

### Token Manipulation Pattern
When adding operations on ANSI strings:
1. Tokenize: `tokens := TokeniseANSIString(input)`
2. Operate on `[][]ANSILineToken` structure
3. Build output: `output := BuildANSIString(tokens, indent)`
4. Never manipulate raw ANSI strings directly

### Test Expectations Format
Write expected results as token arrays, not raw ANSI strings:
```go
expected: [][]convert.ANSILineToken{
    {
        {FG: "", BG: "", T: ""},
        {FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
        {FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: " AAA"},
    },
}
```

## Key Files Reference
- [main.go](main.go): CLI argument parsing (uses `github.com/pborman/getopt/v2`)
- [src/convert/convert.go](src/convert/convert.go): Core tokenization, flip, sanitize logic
- [src/convert/mirror.go](src/convert/mirror.go): Character mirroring lookup tables
- [test/helper.go](test/helper.go): Test utilities and formatting
