# exprgrep

Small CLI that reads JSON Lines (JSONL) from stdin and prints lines that match
an `expr` expression (using https://github.com/expr-lang/expr).

## Install

### Homebrew

```bash
brew tap jtarchie/exprgrep https://github.com/jtarchie/exprgrep
brew install exprgrep
```

### From source

```bash
GOEXPERIMENT=jsonv2 go build -o exprgrep .
```

## Usage

Run (reads JSONL from stdin; expression is first argument):

```bash
cat input.jsonl | ./exprgrep 'age > 30 && active == true'
```

Allow references to fields that may not exist in every record (missing fields
evaluate to `nil` instead of causing an error):

```bash
cat input.jsonl | ./exprgrep --allow-missing-fields 'age != nil && age > 30'
```

Print a specific field (or computed value) instead of the whole original line:

```bash
cat input.jsonl | ./exprgrep --output 'request_id' 'status == 500'
```

This pairs well with a second `rg` pass to find all related log lines:

```bash
rg -z 'status=500' ~/Downloads/*.json.gz \
  | ./exprgrep --output 'request_id' 'action == "create"' \
  | rg -z -f - ~/Downloads/*.json.gz
```

## Exit codes

| Code | Meaning                            |
| ---- | ---------------------------------- |
| 0    | At least one line matched          |
| 1    | No lines matched                   |
| 2    | Error (bad expression, usage, etc) |

## Notes

- The program expects one JSON value per line (JSONL) on stdin.
- The first argument is an `expr` expression evaluated against the parsed JSON
  value.
- If the expression evaluates to a boolean `true`, the original line is printed.
- Non-boolean truthy values (e.g. a non-empty string) also cause the line to be
  printed.
- Lines with invalid JSON are skipped with a warning on stderr.
- This project uses the experimental `encoding/json/v2` API; building with the
  `jsonv2` experiment enabled is required (handled automatically by Homebrew and
  the release builds).

## File

- `main.go` — the CLI implementation
- `main_test.go` — unit tests
