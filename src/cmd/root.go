// Package cmd implements the raci command line. The logic lives in
// raci-manager; commands parse flags, call it and present the results.
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// VersionInfo is the build information printed by `raci version`.
type VersionInfo struct {
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
}

// Licence and project pointers shown by --version, `raci version` and the
// root help. They must match LICENSE and NOTICE, and the URL must resolve.
const (
	licenseName = "Apache 2.0"
	licenseLine = "Apache License 2.0 · Copyright 2026 Consulare and contributors"
	websiteURL  = "https://raci.sh"
)

var versionInfo = VersionInfo{Version: "dev", BuildTime: "unknown", GitCommit: "unknown", GitBranch: "unknown"}

// SetVersionInfo records the build information injected into main.
func SetVersionInfo(v VersionInfo) {
	versionInfo = v
}

// NewRootCmd creates a fresh root command with all subcommands registered.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "raci",
		SilenceErrors: true, // Execute prints the error, sanitized
		SilenceUsage:  true,
		Short:         "RACI for agents",
		Long:          "RACI.sh CLI — who is Responsible, Accountable, Consulted and Informed, per realm. Works offline from your repository.",
		Version:       versionInfo.Version,
	}
	// The first line stays `raci <version>` so scripts that parse it keep working.
	cmd.SetVersionTemplate("raci {{.Version}}\n" + licenseLine + "\n" + websiteURL + "\n")
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newVerifyCmd())
	cmd.SetHelpFunc(customHelp)
	return cmd
}

// Execute runs the root command. Errors are printed once to stderr and the
// process exits with the code carried by an exitError, or 1 otherwise. The
// context is cancelled on SIGINT or SIGTERM.
func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	root := NewRootCmd()
	err := root.ExecuteContext(ctx)
	stop()
	if err != nil {
		_, _ = fmt.Fprintln(root.ErrOrStderr(), printableDiagnostic("Error: "+err.Error()))
		os.Exit(exitCodeOf(err))
	}
}

func customHelp(cmd *cobra.Command, _ []string) {
	w := cmd.OutOrStdout()

	writeln(w, brandStyle.Render(cmd.UseLine()))
	writeln(w)

	desc := cmd.Long
	if strings.TrimSpace(desc) == "" {
		desc = cmd.Short
	}
	if strings.TrimSpace(desc) != "" {
		writeln(w, descStyle.Render("  "+desc))
		writeln(w)
	}

	var visible []*cobra.Command
	for _, c := range cmd.Commands() {
		if !c.Hidden && c.IsAvailableCommand() {
			visible = append(visible, c)
		}
	}
	if len(visible) > 0 {
		writeln(w, brandStyle.Render("Available Commands:"))
		for _, c := range visible {
			writeln(w, cmdStyle.Render(fmt.Sprintf("  %-12s", c.Name()))+descStyle.Render(c.Short))
		}
		writeln(w)
	}

	if local := renderFlags(cmd.LocalFlags()); len(local) > 0 {
		writeln(w, brandStyle.Render("Flags:"))
		writeln(w, strings.Join(local, "\n"))
		writeln(w)
	}
	if inherited := renderFlags(cmd.InheritedFlags()); len(inherited) > 0 {
		writeln(w, brandStyle.Render("Global Flags:"))
		writeln(w, strings.Join(inherited, "\n"))
		writeln(w)
	}

	writeln(w, descStyle.Render(fmt.Sprintf("  Use \"%s [command] --help\" for more information about a command.", cmd.CommandPath())))

	// Only the root help carries the licence pointer; subcommand help stays focused.
	if !cmd.HasParent() {
		writeln(w)
		writeln(w, descStyle.Render(fmt.Sprintf("  Licensed under %s · %s", licenseName, websiteURL)))
	}
}

func writeln(w io.Writer, a ...any) {
	_, _ = fmt.Fprintln(w, a...)
}

func renderFlags(flags *pflag.FlagSet) []string {
	var lines []string
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		short := ""
		if f.Shorthand != "" {
			short = "-" + f.Shorthand + ", "
		}
		lines = append(lines, cmdStyle.Render(fmt.Sprintf("  %s--%-12s", short, f.Name))+descStyle.Render(f.Usage))
	})
	return lines
}
