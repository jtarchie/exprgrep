package main

import (
	"bufio"
	jsonv2 "encoding/json/v2"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/expr-lang/expr"
)

func run(expression string, opts []expr.Option, r io.Reader, w io.Writer) (bool, error) {
	program, err := expr.Compile(expression, opts...)
	if err != nil {
		return false, fmt.Errorf("invalid expression: %w", err)
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	matched := false
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
		out, err := expr.Run(program, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "expr error: %v\n", err)
			continue
		}
		switch v := out.(type) {
		case bool:
			if v {
				matched = true
				fmt.Fprintln(bw, line)
			}
		case nil:
			// no match
		default:
			_ = v
			matched = true
			fmt.Fprintln(bw, line)
		}
	}
	return matched, scanner.Err()
}

func main() {
	log.SetFlags(0)

	allowMissing := flag.Bool("allow-missing-fields", false, "treat missing JSON fields as nil instead of an error")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: exprgrep [--allow-missing-fields] '<expression>'")
		os.Exit(2)
	}
	expression := flag.Arg(0)

	var opts []expr.Option
	if *allowMissing {
		opts = append(opts, expr.AllowUndefinedVariables())
	}

	matched, err := run(expression, opts, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	if !matched {
		os.Exit(1)
	}
}
