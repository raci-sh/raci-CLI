<div align="center">

<img src="https://raci.sh/assets/logo.svg" alt="RACI.sh logo" width="96" height="96">

# RACI.sh CLI

**RACI for agents: who is Responsible, Accountable, Consulted and Informed, per realm.**

[![CI](https://img.shields.io/github/actions/workflow/status/raci-sh/raci-CLI/ci.yml?branch=develop&label=CI&labelColor=161B22)](https://github.com/raci-sh/raci-CLI/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/raci-sh/raci-CLI?labelColor=161B22&color=3FB950)](https://github.com/raci-sh/raci-CLI/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/raci-sh/raci-CLI/develop?labelColor=161B22&color=161B22)](go.mod)
[![License](https://img.shields.io/badge/License-Apache_2.0-161B22?labelColor=161B22)](LICENSE)

[Install](#installation) · [Quick start](#quick-start) · [Commands](#commands) ·
[Governance](#governance) · [Security](#security) · [License](#license)

</div>

`raci` is the [RACI.sh](https://raci.sh) command line, built on the
[raci-manager](https://github.com/raci-sh/raci-manager) engine. It records who is **R**esponsible,
**A**ccountable, **C**onsulted and **I**nformed for activities on resources, per realm.

It works offline from your repository. A realm's RACI lives in a `.raci` file, which is a
materialised view, and its append-only, hash-chained audit log sits next to it. You need nothing
beyond the CLI and GitHub: commit both files and review changes in pull requests. A freemium cloud
service with CI gates and MCP access is planned, and every basic feature will stay available
without it.

RACI.sh realms match [SteerSpec](https://steerspec.dev) realms, but SteerSpec is not required. The
normative specification is in [raci-sh/raci-spec](https://github.com/raci-sh/raci-spec).

**Status.** Early. The [latest release](https://github.com/raci-sh/raci-CLI/releases/latest) ships
`raci verify`, which recomputes the log's hash chain. `raci init` and the replay, rule and change
checks are not implemented yet, and the cloud service is not available.

## Installation

**Homebrew (macOS)**

```bash
brew tap raci-sh/tap
brew trust --cask raci-sh/tap/raci
brew install --cask raci
```

Homebrew can refuse to load a cask from a third-party tap until you trust it, which is what
`brew trust` does.

**GitHub Releases**

Download an archive from [Releases](https://github.com/raci-sh/raci-CLI/releases). Linux and macOS
builds (amd64, arm64) are `.tar.gz` files. Windows builds (amd64, arm64) are `.zip` files with
`raci.exe`, which you can put in any directory on your `PATH`. Each release has a `checksums.txt`.
The macOS binaries are signed and notarized. The Linux and Windows binaries are not code-signed, so
SmartScreen may warn on the first run of the Windows build.

**From source** (Go 1.26+)

```bash
git clone https://github.com/raci-sh/raci-CLI.git
cd raci-CLI
make build    # produces ./raci
```

Use `make build`: `go install` names the binary `src` after the package directory, and leaves out
the version information.

## Quick start

```bash
raci version

# Check the hash chain of the realm's log in the current directory
raci verify --dir . --realm com.example
```

`raci verify` prints a JSON report and exits 1 when it has errors, so it can gate CI.

## Commands

| Command | What it does |
| ------- | ------------ |
| `raci version` | Show version, build time, commit, branch, licence and website. |
| `raci init --realm <id>` | Start a realm's RACI log in this repository. *Not implemented yet.* |
| `raci verify [--dir <dir>] [--realm <id>]` | Recompute the log's hash chain and print a JSON report ([SPEC.md](https://github.com/raci-sh/raci-spec/blob/main/SPEC.md) section 7.1). |
| `raci completion <shell>` | Generate a shell completion script. |

**Files.** `.raci` (the view) and `.raci.jsonl` (the log) sit side by side in the directory that
`--dir` selects. These names are a proposal, since the specification leaves them open.

**Exit codes.** 0 ok, 1 findings, 2 usage or configuration error, 3 not implemented yet.

## Governance

Open source software by Consulare.

Consulare is the copyright holder named in the [NOTICE](NOTICE) file, and the code is licensed under
the [Apache 2.0 licence](LICENSE).

### Development

```bash
make test     # go test -race ./...
make lint     # golangci-lint (gofumpt)
make build    # ./raci
```

RACI logic belongs in [raci-manager](https://github.com/raci-sh/raci-manager): this repository
parses flags, calls the manager and presents results. Commits follow
[Conventional Commits](https://www.conventionalcommits.org), checked locally and in CI.

### Branching and releases

`develop` is the default branch and takes squash merges from feature branches, so a pull request
title must be a conventional commit. Releases reach `main` only through the *Promote to production*
workflow, as merge commits. release-please then tags `vX.Y.Z` and writes the
[changelog](CHANGELOG.md), and goreleaser publishes the binaries and the Homebrew cask.

## Security

Report a vulnerability privately: on this repository, go to **Security**, then **Report a
vulnerability**, or open the [advisory form](https://github.com/raci-sh/raci-CLI/security/advisories/new).
Do not open a public issue, discussion or pull request for a security problem.

## License

Apache 2.0, see [LICENSE](LICENSE).
