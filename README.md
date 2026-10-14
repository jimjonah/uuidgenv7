# uuidgenv7: generate UUID v7 in the macOS terminal

`uuidgenv7` is a small Go CLI that prints UUID version 7 values ([RFC 9562](https://www.rfc-editor.org/rfc/rfc9562)) on macOS. It works like the built-in `uuidgen` command, which only produces random version 4 UUIDs, and installs with Homebrew.

```console
$ uuidgenv7
01a11878-4418-7249-97a3-702d55aec2d1

$ uuidgenv7 -n 3
01a11878-441c-7ae1-8e6b-dcae307145da
01a11878-441c-7aed-8c15-d7444ced5541
01a11878-441c-7af1-89e3-57df594bdb8f
```

## Why UUID v7

A UUIDv7 starts with a millisecond timestamp, so values sort in the order they were created. That makes them a good fit for database primary keys: new rows land at the end of the index instead of at random positions, as they do with version 4 UUIDs. The remaining bits are random, so values are still unique without coordination between machines.

## Install

With Homebrew, on macOS 13 or later (Apple Silicon and Intel):

```sh
brew install jimjonah/tap/uuidgenv7
```

Or with Go, on any platform Go supports:

```sh
go install github.com/jimjonah/uuidgenv7@latest
```

## Usage

```
uuidgenv7 [-n count]
uuidgenv7 -v | --version
uuidgenv7 -h | --help
```

| Option | Meaning |
| --- | --- |
| `-n count` | Number of UUIDs to print, one per line. Default 1. |
| `-v`, `--version` | Print the version and exit. |
| `-h`, `--help` | Print the help and exit. |

Options take one or two dashes, and a value can follow a space or `=`: `-n 5`, `--n 5` and `--n=5` are the same.

Output is lowercase, the canonical form in RFC 9562; macOS `uuidgen` prints uppercase.

UUIDs from a single run are strictly increasing, so the output of `-n` is already sorted. Separate runs in the same millisecond are still unique, but aren't guaranteed to sort in the order they ran.

To save a batch to a file:

```sh
uuidgenv7 -n 1000 > ids.txt
```

The Homebrew install reports its release number, e.g. `0.1.0`. Copies built with Go, including `go install`, report `dev`.

Exit status is 0 on success and 2 for invalid usage, such as `-n 0` or an unknown option.

## License

MIT. See [LICENSE](LICENSE).
