package main

import (
	"os"

	"github.com/conplementag/cops-hq/v2/pkg/error_handling"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
)

// see .goreleaser.yaml ldflags
var version = "local"

func main() {
	defer errorhandler()
	hq := hq.NewQuiet("copsctl", version, "copsctl.log")
	createCommands(hq)

	error_handling.PanicOnAnyError = true

	hq.Run()
}

func errorhandler() {
	// as this is a CLI tool, and not a library with an API, panic is used for most errors that occur,
	// since they are unrecoverable and need some user intervention (or they are genuine panic programming
	// errors)
	if r := recover(); r != nil {
		logrus.Errorf("copsctl --- error occured: %+v\n", r)
		os.Exit(1)
	}
}
