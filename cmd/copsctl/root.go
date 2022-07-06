package main

import (
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"os"
	"strings"

	"github.com/conplementAG/copsctl/internal/common/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copsctl",
	Short: "copsctl - the Conplement AG Kubernetes developer tooling",
	Long: `

                              _   _ 
      ___ ___  _ __  ___  ___| |_| |
     / __/ _ \| '_ \/ __|/ __| __| |
    | (_| (_) | |_) \__ \ (__| |_| |
     \___\___/| .__/|___/\___|\__|_|
              |_|                       by Conplement AG
	
	Version 0.8.2
    `,

	Version: "0.8.2",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(createNamespaceCommand())
	rootCmd.AddCommand(createClusterInfoCommand())
	rootCmd.AddCommand(createConnectCommand())
	rootCmd.AddCommand(createAzureDevopsCommand())

	rootCmd.PersistentFlags().BoolP(flags.Verbose, "v", false, "If set logging will be verbose")
	viper.BindPFlag(flags.Verbose, rootCmd.PersistentFlags().Lookup(flags.Verbose))
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	logging.Initialize()

	logging.Debug("=============================================")
	logging.Debug("===============    copsctl    ===============")
	logging.Debug("=============================================")
	logging.Debug("Command " + strings.Join(os.Args, " "))
}
