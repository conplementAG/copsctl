package main

import "github.com/conplementag/cops-hq/v2/pkg/hq"

func createCommands(hq hq.HQ) {
	createNamespaceCommand(hq)
	createInfoCommands(hq)
	createConnectCommand(hq)
	createAzureDevopsCommand(hq)
}
