package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/raci-sh/raci-CLI/src/internal/testutil"
	"github.com/raci-sh/raci-manager/src/event"
	"github.com/raci-sh/raci-manager/src/eventlog"
)

func TestVersion(t *testing.T) {
	SetVersionInfo(VersionInfo{Version: "1.2.3", GitCommit: "abcdef0123", BuildTime: "now", GitBranch: "develop"})
	out, err := testutil.ExecuteCommand(NewRootCmd(), "version")
	if err != nil {
		t.Fatal(err)
	}
	testutil.AssertContains(t, out, "1.2.3")
	testutil.AssertContains(t, out, "abcdef0")
	testutil.AssertNotContains(t, out, "abcdef0123")
}

func TestHelpListsCommands(t *testing.T) {
	out, err := testutil.ExecuteCommand(NewRootCmd(), "--help")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []string{"init", "verify", "version"} {
		testutil.AssertContains(t, out, c)
	}
	out, _ = testutil.ExecuteCommand(NewRootCmd(), "init", "--help")
	testutil.AssertContains(t, out, "--realm")
}

func TestInit(t *testing.T) {
	_, err := testutil.ExecuteCommand(NewRootCmd(), "init")
	if exitCodeOf(err) != exitUsage {
		t.Errorf("no realm: exit %d (%v)", exitCodeOf(err), err)
	}
	_, err = testutil.ExecuteCommand(NewRootCmd(), "init", "--realm", "com.example")
	if exitCodeOf(err) != exitNotImplemented || !errors.Is(err, errNotImplemented) {
		t.Errorf("realm: exit %d (%v)", exitCodeOf(err), err)
	}
}

func TestVerify(t *testing.T) {
	dir := t.TempDir()
	out, err := testutil.ExecuteCommand(NewRootCmd(), "verify", "--dir", dir, "--realm", "sh.raci")
	if err != nil {
		t.Fatalf("empty log: %v\n%s", err, out)
	}
	testutil.AssertContains(t, out, `"findings": []`)

	log := eventlog.Open(filepath.Join(dir, eventlog.FileName))
	if err := log.Append(&event.Event{ID: "e1", Type: event.TypeCommit, Payload: []byte(`{}`)}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(log.Path())
	if err := os.WriteFile(log.Path(), append([]byte(`{"seq":1,"prev_hash":"x"}`+"\n"), data...), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err = testutil.ExecuteCommand(NewRootCmd(), "verify", "--dir", dir)
	if exitCodeOf(err) != exitFindings {
		t.Errorf("tampered: exit %d (%v)", exitCodeOf(err), err)
	}
	testutil.AssertContains(t, out, "hash-mismatch")

	if err := os.WriteFile(log.Path(), []byte("garbage\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := testutil.ExecuteCommand(NewRootCmd(), "verify", "--dir", dir); exitCodeOf(err) != exitUsage {
		t.Errorf("unreadable: exit %d (%v)", exitCodeOf(err), err)
	}
}

func TestExitHelpers(t *testing.T) {
	if exitCodeOf(nil) != exitOK || withExitCode(exitUsage, nil) != nil || exitCodeOf(errors.New("x")) != exitFindings {
		t.Error("exit helpers")
	}
	if got := printableDiagnostic("a\x1b[31mb\n"); got != "a[31mb\n" {
		t.Errorf("printableDiagnostic = %q", got)
	}
}
