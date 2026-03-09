# exprgrep

Small CLI that reads JSON Lines (JSONL) from stdin and prints lines that match
an `expr` expression (using https://github.com/expr-lang/expr).

## Usage

Build:

```bash
GOEXPERIMENT=jsonv2 go build -o exprgrep .
```

Run (reads JSONL from stdin; expression is first argument):

```bash
cat input.jsonl | ./exprgrep 'age > 30 && active == true'
```

## Notes

- The program expects one JSON value per line (JSONL) on stdin.
- The first argument is an `expr` expression evaluated against the parsed JSON
  value.
- If the expression evaluates to a boolean `true`, the original line is printed.
- This project was tested using the experimental `encoding/json/v2` API;
  building with the `jsonv2` experiment enabled is recommended to reproduce the
  same behavior:

```bash
GOEXPERIMENT=jsonv2 go build -o exprgrep .
```

Optional: fetch the `expr` module (the build will do this automatically):

```bash
go get github.com/expr-lang/expr@latest
```

## File

- `main.go` — the CLI implementation
