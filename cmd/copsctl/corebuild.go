package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/common"
	"github.com/conplementAG/copsctl/internal/corebuild"
	"github.com/conplementag/cops-hq/v2/pkg/cli"
	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func createCoreBuildCommands(hq hq.HQ) {
	coreBuildCmdGroup := hq.GetCli().AddBaseCommand("build",
		"Command group for core build tasks",
		"Use this command to manipulate core build infrastructure", nil)

	createCoreBuildCreateCommand(coreBuildCmdGroup, hq)
	createCoreBuildDestroyCommand(coreBuildCmdGroup, hq)
}

func createCoreBuildCreateCommand(cmd cli.Command, hq hq.HQ) {
	createCoreBuildCmd := cmd.AddCommand("create", "Create core build infrastructure",
		"Use this command to create core build infrastructure including azure resources and azure devops configurations",
		func() {
			orchestrator, err := corebuild.New(hq)
			common.FatalOnError(err)

			orchestrator.CreateInfrastructure()
		})

	addConfigFileParam(createCoreBuildCmd)
}

func createCoreBuildDestroyCommand(cmd cli.Command, hq hq.HQ) {
	createCoreBuildCmd := cmd.AddCommand("destroy", "Destroy core build infrastructure",
		"Use this command to destroy core build infrastructure including azure resources and azure devops configurations",
		func() {
			orchestrator, err := corebuild.New(hq)
			common.FatalOnError(err)

			orchestrator.DestroyInfrastructure()
		})

	addConfigFileParam(createCoreBuildCmd)
}

func addConfigFileParam(cmd cli.Command) {
	cmd.AddPersistentParameterString(flags.ConfigFile, "", true, "f",
		"Yaml config for createing core build infrastructure")
	cmd.AddPersistentParameterString(flags.SopsConfigFile, "", true, "c",
		"Configuration file path for sops configuration file. If not given sops config is expected next to config file.")
}
