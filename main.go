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

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("uuidgenv7", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: uuidgenv7 [-n count] [-version]")
		fs.PrintDefaults()
	}
	n := fs.Int("n", 1, "number of UUIDs to generate")
	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "uuidgenv7: unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	if *showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if *n < 1 {
		fmt.Fprintln(stderr, "uuidgenv7: -n must be at least 1")
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
