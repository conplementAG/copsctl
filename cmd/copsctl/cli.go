package main

import "github.com/conplementag/cops-hq/pkg/hq"

func createCommands(hq hq.HQ) {
	createNamespaceCommand(hq)
	createClusterInfoCommand(hq)
	createConnectCommand(hq)
	createAzureDevopsCommand(hq)
}
