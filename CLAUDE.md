# CLAUDE.md

Guidance for AI coding agents working in this repository.

## What this repo is

`raci-CLI` builds `raci`, the RACI.sh command line: RACI for agents (who is Responsible, Accountable, Consulted and Informed for activities on resources, per realm). It is the single user-facing tool. Locally it works offline from a repository, with a `.raci` view plus its JSONL log. Later it will also drive the cloud mode (`raci-cloud`).

The logic lives in [raci-sh/raci-manager](https://github.com/raci-sh/raci-manager), imported as a Go module. This repo parses flags, calls the manager and presents results. Do not put RACI logic here: add it to the manager, release it, then bump `go.mod`.

The normative source is the specification in [raci-sh/raci-spec](https://github.com/raci-sh/raci-spec) (`SPEC.md`).

Everything basic must work without the cloud, using only the CLI and GitHub.

<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:1105d646 -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/core-concepts/sync-concepts.md for details and anti-patterns.

## Agent Context Profiles

The managed Beads block is task-tracking guidance, not permission to override repository, user, or orchestrator instructions.

- **Conservative (default)**: Use `bd` for task tracking. Do not run git commits, git pushes, or Dolt remote sync unless explicitly asked. At handoff, report changed files, validation, and suggested next commands.
- **Minimal**: Keep tool instruction files as pointers to `bd prime`; use the same conservative git policy unless active instructions say otherwise.
- **Team-maintainer**: Only when the repository explicitly opts in, agents may close beads, run quality gates, commit, and push as part of session close. A current "do not commit" or "do not push" instruction still wins.

## Session Completion

This protocol applies when ending a Beads implementation workflow. It is subordinate to explicit user, repository, and orchestrator instructions.

1. **File issues for remaining work** - Create beads for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **Handle git/sync by active profile**:
   ```bash
   # Conservative/minimal/default: report status and proposed commands; wait for approval.
   git status

   # Team-maintainer opt-in only, unless current instructions forbid it:
   git pull --rebase
   git push
   git status
   ```
5. **Hand off** - Summarize changes, validation, issue status, and any blocked sync/commit/push step

**Critical rules:**
- Explicit user or orchestrator instructions override this Beads block.
- Do not commit or push without clear authority from the active profile or the current user request.
- If a required sync or push is blocked, stop and report the exact command and error.
<!-- END BEADS INTEGRATION -->

## Build & Test

```sh
make build   # ./raci (main package is ./src, so always build with -o raci)
make test    # go test -race ./...
make lint    # golangci-lint v2 (gofumpt)
make fmt
```

## Architecture

- `src/main.go`: build info vars (`Version`, `GitCommit`, `BuildTime`, `GitBranch`) injected with `-X main.…`, then `cmd.Execute()`.
- `src/cmd/`: one file per command (`version.go`, `init.go`, `verify.go`), plus:
  - `root.go`: `NewRootCmd`, a signal-aware `Execute`, and the lipgloss help.
  - `exit.go`: exit codes. 0 ok, 1 findings, 2 usage/config, 3 not implemented.
  - `styles.go`.
- `src/internal/ui`: colours. `src/internal/testutil`: `ExecuteCommand` and `AssertContains`.
- Commands return errors wrapped with `withExitCode`; only `Execute` prints errors and exits.
- Offline layout, a scaffold proposal since the spec leaves names open: `.raci` (view) and `.raci.jsonl` (log) side by side; `--dir` selects the directory.

## Conventions

- Conventional Commits. `.beads/hooks/commit-msg` checks locally; `lint-commits` and `pr-title` check in CI.
- Each command has tests in `src/cmd/*_test.go` using `testutil.ExecuteCommand`.
- Apache 2.0. Copyright: Consulare and contributors (`LICENSE`, `NOTICE`).

## Branching and releases

- The default branch is **`develop`**. Feature branch → PR into `develop`, **squash only**; the PR title becomes the commit, so it must be a conventional commit.
- `develop` → `main` happens only through the **Promote to production** workflow (`gh workflow run promote.yml --ref main`). It uses a **merge commit only**, then merges release-please's release PR and the sync-back PR.
- release-please tags `vX.Y.Z` on `main` and writes `CHANGELOG.md`. The tag runs goreleaser (`release.yml`).
- `sync-release-files.yml` carries `CHANGELOG.md` and the manifest back to `develop`. Never hand-edit them.
- Rulesets `merge-method-develop` (squash) and `merge-method-main` (merge) enforce the methods.
- Agents open PRs; they never merge them.
