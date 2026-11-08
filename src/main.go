// Command raci is the RACI.sh CLI: RACI for agents, offline and repo based.
package main

import "github.com/raci-sh/raci-CLI/src/cmd"

// Build information, set through -ldflags by the Makefile and goreleaser.
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

func main() {
	cmd.SetVersionInfo(cmd.VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
	})
	cmd.Execute()
}
