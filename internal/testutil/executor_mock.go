// Package testutil provides shared test doubles for the copsctl test suite.
package testutil

import (
	"os/exec"

	"github.com/conplementag/cops-hq/v2/pkg/commands"
)

// ExecutorMock is a minimal commands.Executor test double. All string-returning Execute*
// methods route through ExecuteFunc (or the static Output/Err fallback); the *Cmd/TTY/confirm
// methods are not needed by the code under test and panic if unexpectedly called.
type ExecutorMock struct {
	// ExecuteFunc, if set, handles every Execute* call and lets a test script per-call behaviour.
	ExecuteFunc func(command string) (string, error)
	// Output and Err are returned when ExecuteFunc is nil.
	Output string
	Err    error
}

// compile-time check that the mock satisfies the interface
var _ commands.Executor = (*ExecutorMock)(nil)

func (m *ExecutorMock) run(command string) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command)
	}
	return m.Output, m.Err
}

func (m *ExecutorMock) Execute(command string) (string, error) { return m.run(command) }
func (m *ExecutorMock) ExecuteWithProgressInfo(command string) (string, error) {
	return m.run(command)
}
func (m *ExecutorMock) ExecuteSilent(command string) (string, error) { return m.run(command) }
func (m *ExecutorMock) ExecuteLoud(command string) (string, error)   { return m.run(command) }

func (m *ExecutorMock) ExecuteCmd(cmd *exec.Cmd) (string, error) { panic("not implemented") }
func (m *ExecutorMock) ExecuteCmdWithProgressInfo(cmd *exec.Cmd) (string, error) {
	panic("not implemented")
}
func (m *ExecutorMock) ExecuteCmdSilent(cmd *exec.Cmd) (string, error) { panic("not implemented") }
func (m *ExecutorMock) ExecuteTTY(command string) error               { panic("not implemented") }
func (m *ExecutorMock) ExecuteCmdTTY(cmd *exec.Cmd) error             { panic("not implemented") }
func (m *ExecutorMock) AskUserToConfirm(displayMessage string) bool   { panic("not implemented") }
func (m *ExecutorMock) AskUserToConfirmWithKeyword(displayMessage string, keyword string) bool {
	panic("not implemented")
}
