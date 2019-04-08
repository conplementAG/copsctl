package commands

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
)

// ExecuteCommand executes the command and respects the current verbosity level
func ExecuteCommand(command *exec.Cmd) string {
	return execute(command, viper.GetBool("verbose"))
}

// ExecuteCommandVerbose executes the given command, pipes to stdout and stderr
// and returns the stdout for further processing, no matter what the verbosity is
func ExecuteCommandVerbose(command *exec.Cmd) string {
	return execute(command, true)
}

// ExecuteCommandLongRunning runs the given command and starts
// a async spinner to tell the user to be patient
func ExecuteCommandLongRunning(command *exec.Cmd) string {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Prefix = "Please wait "
	spinner.Color("green", "bold")
	spinner.Start()

	result := execute(command, viper.GetBool("verbose"))

	spinner.Stop()

	return result
}

func execute(command *exec.Cmd, isVerbose bool) string {
	command.Stdin = os.Stdin

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()

	var errStdout, errStderr error

	// all the error should be piped to Stderr, regardless of verbosity,
	// but for the strout only if verbose is requested
	errorWriter := os.Stderr
	stdWriter := ioutil.Discard

	if isVerbose {
		stdWriter = os.Stdout
	}

	stdout := io.MultiWriter(stdWriter, &stdoutBuf)
	stderr := io.MultiWriter(errorWriter, &stderrBuf)

	if isVerbose {
		log.Printf("Executing: %s %s\n", command.Path, strings.Join(command.Args[1:], " "))
	}

	err := command.Start()

	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	var multiWritingSteps sync.WaitGroup
	multiWritingSteps.Add(2)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		multiWritingSteps.Done()
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
		multiWritingSteps.Done()
	}()

	errCmd := command.Wait()
	multiWritingSteps.Wait()

	handleError(errCmd, errStdout, errStderr, stdoutBuf)

	outStr := string(stdoutBuf.Bytes())
	return outStr
}

func handleError(errCmd error, errStdout error, errStderr error, stdout bytes.Buffer) {
	if errCmd != nil {
		// if the stdout has some info, print it out as well because it might be helpful
		// to understand the error
		if stdout.Len() > 0 {
			log.Println(string(stdout.Bytes()))
		}

		panic(errCmd)
	}
}

// ExecuteCommandForUserInput starts the given command and waits either for command completion or for a user input (ENTER).
// If there is a user input a Kill-Signal will be sent to the command-process and the execution is continued.
func ExecuteCommandForUserInput(cmd *exec.Cmd) {
	isVerbose := viper.GetBool("verbose")

	if isVerbose {
		log.Printf("Executing async: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	cmd.Start()

	done := make(chan error)

	go func() { done <- cmd.Wait() }()
	userInput := make(chan int, 1)

	go func() {
		log.Println("[Press Enter key to continue] ")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadRune()
		userInput <- 1
	}()

	select {
	case <-userInput:
		cmd.Process.Signal(os.Kill)
		if isVerbose {
			log.Println("User completes the command")
		}
	case err := <-done:
		if isVerbose {
			log.Println("Command finished by itself")
		}
		if err != nil {
			panic(err)
		}
	}
}

func ExecutesSuccessfully(command *exec.Cmd) bool {
	if viper.GetBool("verbose") {
		log.Printf("Executing: %s %s\n", command.Path, strings.Join(command.Args[1:], " "))
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	return err == nil
}
