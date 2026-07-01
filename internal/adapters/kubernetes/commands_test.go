package kubernetes

import (
	"errors"
	"testing"

	"github.com/conplementAG/copsctl/internal/testutil"
	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/stretchr/testify/assert"
)

func Test_CanIGetPods_result(t *testing.T) {
	data := []struct {
		testName string
		output   string
		err      error
		expected bool
	}{
		{testName: "yes grants access", output: "yes\n", err: nil, expected: true},
		{testName: "no denies access", output: "no\n", err: nil, expected: false},
		{testName: "error denies access", output: "no\n", err: errors.New("exit status 1"), expected: false},
		{testName: "unexpected output denies access", output: "maybe", err: nil, expected: false},
	}

	for _, test := range data {
		t.Run(test.testName, func(t *testing.T) {
			executor := &testutil.ExecutorMock{Output: test.output, Err: test.err}

			assert.Equal(t, test.expected, CanIGetPods(executor, "some-namespace"))
		})
	}
}

func Test_CanIGetPods_restoresPanicOnAnyError(t *testing.T) {
	data := []struct {
		testName string
		initial  bool
	}{
		{testName: "restores true", initial: true},
		{testName: "restores false", initial: false},
	}

	for _, test := range data {
		t.Run(test.testName, func(t *testing.T) {
			original := error_handling.PanicOnAnyError
			defer func() { error_handling.PanicOnAnyError = original }()
			error_handling.PanicOnAnyError = test.initial

			// an error would normally trip the global panic behaviour - CanIGetPods must tolerate it
			executor := &testutil.ExecutorMock{Output: "no\n", Err: errors.New("exit status 1")}

			CanIGetPods(executor, "some-namespace")

			assert.Equal(t, test.initial, error_handling.PanicOnAnyError, "PanicOnAnyError must be restored to its previous value")
		})
	}
}

func Test_CanIGetPods_disablesPanicOnAnyErrorDuringCommand(t *testing.T) {
	original := error_handling.PanicOnAnyError
	defer func() { error_handling.PanicOnAnyError = original }()
	error_handling.PanicOnAnyError = true

	var flagDuringCommand bool
	executor := &testutil.ExecutorMock{
		ExecuteFunc: func(command string) (string, error) {
			flagDuringCommand = error_handling.PanicOnAnyError
			return "yes\n", nil
		},
	}

	CanIGetPods(executor, "some-namespace")

	assert.False(t, flagDuringCommand, "PanicOnAnyError must be disabled while the command runs")
}
