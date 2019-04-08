package commands

import (
	"fmt"
	"os/exec"
	"strings"
)

const SPACE_WITHIN_QUOTE = "{SPACE_WITHIN_QUOTE}"
const DOUBLE_QUOTE = "{DOUBLE_QUOTE}"

// Create creates a command from a single string
// This allows you to pass parameters which include spaces as commands.
// You just need to add "double-quotes" around the parameter and it will be treated as one parameter and not be splitted by whitespace.
func Create(plainCommand string) *exec.Cmd {
	/*
		Example
		     plainCommand   : az role assignment create --role "Network Contributer" --assignee ABC --scope abc
		     escapedCommand : az role assignment create --role "Network{SPACE_WITHIN_QUOTE}Contributer" --assignee ABC --scope abc
		     commandParts   : ["az", "role", "assignment", "create", "--role", "Network Contributer", "--assignee", "ABC", "--scope", "abc"]
	*/
	escapedCommand := markQuotedSpaces(plainCommand)
	commandParts := handleSpacesInQuotes(strings.Fields(escapedCommand))

	var cmd *exec.Cmd

	if len(commandParts) > 1 {
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	} else {
		cmd = exec.Command(commandParts[0])
	}

	return cmd
}

// Createf creates a command from a formatted input string
func Createf(format string, args ...interface{}) *exec.Cmd {
	actualCommand := fmt.Sprintf(format, args...)
	return Create(actualCommand)
}

func markQuotedSpaces(plainCommand string) string {
	escapeMode := false
	var escapedCommand strings.Builder
	enteredEscapeModeInThisIteration := false

	for _, char := range plainCommand {

		// Enter Escape Mode when " occurrs
		if !escapeMode && string(char) == "\"" {
			escapeMode = true
			enteredEscapeModeInThisIteration = true
		}

		// Handle spaces when in Escape Mode
		if escapeMode && string(char) == " " {
			escapedCommand.WriteString(SPACE_WITHIN_QUOTE)
		} else if escapeMode && string(char) == "\"" {
			escapedCommand.WriteString(DOUBLE_QUOTE)
		} else {
			escapedCommand.WriteRune(char)
		}

		// Exit Escape Mode when " occurrs
		if !enteredEscapeModeInThisIteration && escapeMode && string(char) == "\"" {
			escapeMode = false
		}

		enteredEscapeModeInThisIteration = false
	}
	return escapedCommand.String()
}

func handleSpacesInQuotes(parts []string) []string {
	handledParts := make([]string, len(parts))

	for pos, part := range parts {
		newPartWithSpaces := strings.Replace(part, SPACE_WITHIN_QUOTE, " ", -1)
		newPartWithoutQuotes := strings.Replace(newPartWithSpaces, DOUBLE_QUOTE, "", -1)
		handledParts[pos] = newPartWithoutQuotes
	}

	return handledParts
}
