package main

import (
	"bufio"
	jsonv2 "encoding/json/v2"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/expr-lang/expr"
)

func main() {
    log.SetFlags(0)
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: exprgrep '<expression>'")
        os.Exit(2)
    }
    expression := os.Args[1]

    scanner := bufio.NewScanner(os.Stdin)
    w := bufio.NewWriter(os.Stdout)
    defer w.Flush()

    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            continue
        }
        var data interface{}
        if err := jsonv2.Unmarshal([]byte(line), &data); err != nil {
            fmt.Fprintf(os.Stderr, "invalid json: %v\n", err)
            continue
        }
        out, err := expr.Eval(expression, data)
        if err != nil {
            fmt.Fprintf(os.Stderr, "expr error: %v\n", err)
            continue
        }
        match := false
        switch v := out.(type) {
        case bool:
            match = v
        case nil:
            match = false
        default:
            match = true
            _ = v
        }
        if match {
            fmt.Fprintln(w, line)
        }
    }
    if err := scanner.Err(); err != nil && err != io.EOF {
        log.Fatal(err)
    }
}
