# raci

The [RACI.sh](https://raci.sh) CLI: RACI for agents. It records who is **R**esponsible, **A**ccountable, **C**onsulted and **I**nformed for activities on resources, per realm.

`raci` works offline from your repository. A realm's RACI lives in a `.raci` file, which is a materialised view. Its append-only, hash-chained audit log sits next to it. You need nothing beyond the CLI and GitHub: commit both files and review changes in pull requests. A freemium cloud service with CI gates and MCP access will come later, but every basic feature stays available without it.

RACI.sh is a companion to [SteerSpec](https://github.com/SteerSpec), and its realms match SteerSpec realms. It does not require SteerSpec. The specification is in [raci-sh/raci-spec](https://github.com/raci-sh/raci-spec). The engine is [raci-sh/raci-manager](https://github.com/raci-sh/raci-manager).

> Status: scaffold. `raci verify` checks the log's hash chain. `raci init` and everything else are still to be built.

## Install

From a release: download the archive for your platform from [Releases](https://github.com/raci-sh/raci-CLI/releases).

From source:

```sh
git clone https://github.com/raci-sh/raci-CLI && cd raci-CLI
make build   # produces ./raci
```

## Usage

```sh
raci version
raci init --realm com.example   # not implemented yet
raci verify [--dir .]           # JSON report; exit 1 on errors
```

Exit codes: 0 ok, 1 findings, 2 usage or configuration error, 3 not implemented yet.

## Development

```sh
make test lint build
```

`CLAUDE.md` covers the branching model. `develop` is the default branch and takes squash merges. Releases reach `main` through the *Promote to production* workflow, as merge commits, and release-please tags them.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
