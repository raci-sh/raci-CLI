package cmd

import (
	"errors"
	"strings"
	"unicode"
)

// Process exit codes. Commands that do not wrap their errors in exitError
// exit with 1.
const (
	exitOK             = 0
	exitFindings       = 1 // verification found errors (or any plain error)
	exitUsage          = 2 // bad flags, arguments or configuration
	exitNotImplemented = 3 // the command or mode is not built yet
)

// errNotImplemented marks commands that are scaffolded but not built yet.
var errNotImplemented = errors.New("not implemented yet")

// exitError carries a process exit code alongside the underlying error.
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string { return e.err.Error() }

func (e *exitError) Unwrap() error { return e.err }

// withExitCode wraps err so that Execute exits with code. A nil err is
// returned unchanged.
func withExitCode(code int, err error) error {
	if err == nil {
		return nil
	}
	return &exitError{code: code, err: err}
}

// exitCodeOf returns the exit code err should produce: the wrapped code for
// an exitError, 1 for any other error and 0 for nil.
func exitCodeOf(err error) int {
	if err == nil {
		return exitOK
	}
	var ee *exitError
	if errors.As(err, &ee) {
		return ee.code
	}
	return exitFindings
}

// printableDiagnostic drops control characters, other than newline and tab,
// so that an error built from file content cannot drive the terminal.
func printableDiagnostic(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			return r
		}
		return -1
	}, s)
}
