package main

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var v7Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestNextFormatAndTimestamp(t *testing.T) {
	before := time.Now().UnixMilli()
	id, err := newGenerator().next()
	after := time.Now().UnixMilli()
	if err != nil {
		t.Fatal(err)
	}
	if !v7Pattern.MatchString(id) {
		t.Fatalf("%q is not a canonical UUIDv7", id)
	}
	ms, err := strconv.ParseInt(id[0:8]+id[9:13], 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	if ms < before || ms > after {
		t.Fatalf("timestamp %d outside [%d, %d]", ms, before, after)
	}
}

func TestNextStrictlyIncreasing(t *testing.T) {
	g := newGenerator()
	prev := ""
	for i := 0; i < 100000; i++ {
		id, err := g.next()
		if err != nil {
			t.Fatal(err)
		}
		if id <= prev {
			t.Fatalf("id %d not increasing: %s then %s", i, prev, id)
		}
		prev = id
	}
}

func TestNextMonotonicWhenClockStalls(t *testing.T) {
	frozen := time.Unix(1700000000, 0)
	g := &generator{now: func() time.Time { return frozen }}
	prev := ""
	// More than 4096 so the 12-bit fraction carries into the milliseconds.
	for i := 0; i < 10000; i++ {
		id, err := g.next()
		if err != nil {
			t.Fatal(err)
		}
		if !v7Pattern.MatchString(id) {
			t.Fatalf("%q is not a canonical UUIDv7", id)
		}
		if id <= prev {
			t.Fatalf("id %d not increasing: %s then %s", i, prev, id)
		}
		prev = id
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantCode  int
		wantLines int
	}{
		{"default", nil, 0, 1},
		{"n", []string{"-n", "5"}, 0, 5},
		{"n equals form", []string{"-n=3"}, 0, 3},
		{"n double dash", []string{"--n", "2"}, 0, 2},
		{"zero", []string{"-n", "0"}, 2, 0},
		{"negative", []string{"-n", "-1"}, 2, 0},
		{"not a number", []string{"-n", "abc"}, 2, 0},
		{"stray argument", []string{"5"}, 2, 0},
		{"unknown flag", []string{"-x"}, 2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != tt.wantCode {
				t.Fatalf("exit code %d, want %d (stderr: %s)", code, tt.wantCode, stderr.String())
			}
			var lines []string
			if stdout.Len() > 0 {
				lines = strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			}
			if len(lines) != tt.wantLines {
				t.Fatalf("got %d lines, want %d: %q", len(lines), tt.wantLines, stdout.String())
			}
			for _, l := range lines {
				if !v7Pattern.MatchString(l) {
					t.Errorf("%q is not a canonical UUIDv7", l)
				}
			}
			if tt.wantCode != 0 && stderr.Len() == 0 {
				t.Error("expected an error message on stderr")
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	for _, arg := range []string{"-h", "--help", "-help"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run([]string{arg}, &stdout, &stderr); code != 0 {
				t.Fatalf("exit code %d", code)
			}
			if stdout.String() != usage {
				t.Errorf("stdout is %q, want the usage text", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Errorf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestRunErrorShowsHint(t *testing.T) {
	var stdout, stderr bytes.Buffer
	run([]string{"-x"}, &stdout, &stderr)
	if !strings.Contains(stderr.String(), "flag provided but not defined: -x") {
		t.Errorf("stderr missing the error: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), helpHint) {
		t.Errorf("stderr missing the help hint: %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("unexpected stdout: %q", stdout.String())
	}
}

func TestRunVersion(t *testing.T) {
	for _, arg := range []string{"-version", "--version", "-v", "--v"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run([]string{arg}, &stdout, &stderr); code != 0 {
				t.Fatalf("exit code %d", code)
			}
			if got := strings.TrimSpace(stdout.String()); got != version {
				t.Fatalf("got %q, want %q", got, version)
			}
		})
	}
}
