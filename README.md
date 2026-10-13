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

UUIDs from a single invocation are strictly increasing, so `-n` output is already sorted.

## Install

```sh
brew install jimjonah/tap/uuidgenv7
```

Or with Go:

```sh
go install github.com/jimjonah/uuidgenv7@latest
```

## Usage

| Flag | Meaning |
| --- | --- |
| `-n count` | Number of UUIDs to print, one per line. Default 1. |
| `-version` | Print the version and exit. |

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


