package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/connection"
	"github.com/conplementag/cops-hq/pkg/hq"
)

func createConnectCommand(hq hq.HQ) {
	connectCmdGroup := hq.GetCli().AddBaseCommand("connect", "Command for managing the connection to k8s clusters",
		"Use this command to manage the connection to a k8s cluster.", func() {
			connection.New(hq).Connect()
		})

	connectCmdGroup.AddPersistentParameterString(flags.EnvironmentTag, "", true, "e",
		"The environment tag of the environment you want to connect to")
	connectCmdGroup.AddPersistentParameterString(flags.ConnectionString, "", true, "c",
		"The connection string for the environment.")
	connectCmdGroup.AddPersistentParameterBool(flags.TechnicalAccount, false, false, "t",
		"Use to switch to non-interactive technical account mode. "+
			"Make sure that the connection string you provide is valid for a technical account!")
	connectCmdGroup.AddPersistentParameterBool(flags.Secondary, false, false, "s",
		"Connect to the secondary instead of the primary cluster.")
	connectCmdGroup.AddPersistentParameterBool(flags.AutoApprove, false, false, "a",
		"Auto approve any confirmation prompts.")
	connectCmdGroup.AddPersistentParameterBool(flags.PrintToStdout, false, false, "p",
		"Print the config to stdout instead of overwriting $HOME/.kube/config. Useful for developers")
	connectCmdGroup.AddPersistentParameterBool(flags.PrintToStdoutSilenceEverythingElse, false, false, "q",
		"Similar to print-to-stdout, but silences all other logging outputs. Useful for automation.")
}
