package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/conplementAG/copsctl/pkg/common/logging"
)

// ExecuteCommandLongRunning runs the given command and starts
// a async spinner to tell the user to be patient
func ExecuteCommandLongRunning(command *exec.Cmd) (string, error) {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Prefix = "Please wait "
	spinner.Color("green", "bold")
	spinner.Start()

	result, err := ExecuteCommand(command)

	spinner.Stop()

	return result, err
}

func ExecuteCommand(command *exec.Cmd) (string, error) {
	command.Stdin = os.Stdin

	commandStdoutPipe, _ := command.StdoutPipe()
	commandstderrPipe, _ := command.StderrPipe()

	var stdoutBuffer, stderrBuffer bytes.Buffer
	stdoutMultiwriter := io.MultiWriter(newDebugLogWriter(), &stdoutBuffer)
	stderrMultiwriter := io.MultiWriter(newDebugLogWriter(), &stderrBuffer)

	logging.Debugf("Executing: %s %s", command.Path, strings.Join(command.Args[1:], " "))

	err := command.Start()

	if err != nil {
		logging.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var multiWritingSteps sync.WaitGroup
	multiWritingSteps.Add(2)

	go func() {
		io.Copy(stdoutMultiwriter, commandStdoutPipe)
		multiWritingSteps.Done()
	}()

	go func() {
		io.Copy(stderrMultiwriter, commandstderrPipe)
		multiWritingSteps.Done()
	}()

	commandReturnValue := command.Wait()
	multiWritingSteps.Wait()

	outStr := string(stdoutBuffer.Bytes())
	return outStr, commandReturnValue
}

type debugLogWriter struct{}

func newDebugLogWriter() *debugLogWriter {
	return &debugLogWriter{}
}

func (w *debugLogWriter) Write(p []byte) (int, error) {
	logging.Debug(string(p))

	return len(p), nil
}
