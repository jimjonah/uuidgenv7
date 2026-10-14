// Command uuidgenv7 prints one or more UUIDv7 values, one per line.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// version is set at release time via -ldflags "-X main.version=...".
var version = "dev"

const usage = `uuidgenv7 prints UUIDv7 values (RFC 9562), one per line.

Usage:
  uuidgenv7 [-n count]
  uuidgenv7 -v | --version
  uuidgenv7 -h | --help

Options:
  -n count       number of UUIDs to print (default 1)
  -v, --version  print the version and exit
  -h, --help     print this help and exit

Options take one or two dashes, and a value can follow a space or "=":
-n 5, --n 5 and --n=5 are the same.

UUIDs from a single run are strictly increasing, so the output of -n is
already sorted.

Examples:
  uuidgenv7
  uuidgenv7 -n 10
  uuidgenv7 -n 1000 > ids.txt
`

const helpHint = "Run 'uuidgenv7 --help' for usage."

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("uuidgenv7", flag.ContinueOnError)
	fs.SetOutput(stderr)
	// Help is printed below rather than by the flag package, so that an
	// explicit -h goes to stdout and usage errors go to stderr.
	fs.Usage = func() {}
	n := fs.Int("n", 1, "")
	var showVersion bool
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&showVersion, "v", false, "")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			fmt.Fprint(stdout, usage)
			return 0
		}
		// The flag package has already printed the error itself.
		fmt.Fprintln(stderr, helpHint)
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "uuidgenv7: unexpected argument %q\n%s\n", fs.Arg(0), helpHint)
		return 2
	}
	if showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if *n < 1 {
		fmt.Fprintf(stderr, "uuidgenv7: -n must be at least 1\n%s\n", helpHint)
		return 2
	}

	w := bufio.NewWriter(stdout)
	g := newGenerator()
	for i := 0; i < *n; i++ {
		id, err := g.next()
		if err != nil {
			fmt.Fprintln(stderr, "uuidgenv7:", err)
			return 1
		}
		w.WriteString(id)
		w.WriteByte('\n')
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintln(stderr, "uuidgenv7:", err)
		return 1
	}
	return 0
}
