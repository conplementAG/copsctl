package main

import "github.com/conplementag/cops-hq/v2/pkg/hq"

func createCommands(hq hq.HQ) {
	createNamespaceCommand(hq)
	createClusterInfoCommand(hq)
	createConnectCommand(hq)
	createAzureDevopsCommand(hq)
}
