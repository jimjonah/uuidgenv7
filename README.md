# uuidgenv7

Print UUIDv7 values ([RFC 9562](https://www.rfc-editor.org/rfc/rfc9562)) from the terminal.

```console
$ uuidgenv7
01a11878-4418-7249-97a3-702d55aec2d1

$ uuidgenv7 -n 3
01a11878-441c-7ae1-8e6b-dcae307145da
01a11878-441c-7aed-8c15-d7444ced5541
01a11878-441c-7af1-89e3-57df594bdb8f
```

## Install

```sh
brew install jimjonah/tap/uuidgenv7
```

Or with Go:

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

UUIDs from a single run are strictly increasing, so the output of `-n` is already sorted. Separate runs in the same millisecond are still unique, but aren't guaranteed to sort in the order they ran.

To save a batch to a file:

```sh
uuidgenv7 -n 1000 > ids.txt
```

A local build reports its version as `dev`; released builds report the release number, e.g. `0.1.0`.

Exit status is 0 on success and 2 for invalid usage, such as `-n 0` or an unknown option.

## Releasing

Pushing a `v*` tag runs `.github/workflows/release.yml`, which uses [GoReleaser](https://goreleaser.com) to upload macOS binaries (Intel and Apple Silicon) to a GitHub release and push an updated cask to `jimjonah/homebrew-tap`.

One-time setup:

1. Create two public GitHub repos: `jimjonah/uuidgenv7` and `jimjonah/homebrew-tap`. The tap can start empty.
2. Create a fine-grained personal access token with **Contents: read and write** on `homebrew-tap`.
3. In the `uuidgenv7` repo, go to Settings → Secrets and variables → Actions and add the token as `HOMEBREW_TAP_GITHUB_TOKEN`.

To cut a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

To check the config locally without publishing anything:

```sh
brew install goreleaser
goreleaser check
goreleaser release --snapshot --clean
```


