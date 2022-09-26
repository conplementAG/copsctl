package connection

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/conplementag/cops-hq/pkg/hq"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
)

type Orchestrator struct {
	hq       hq.HQ
	executor commands.Executor
}

func New(hq hq.HQ) *Orchestrator {
	return &Orchestrator{
		hq:       hq,
		executor: hq.GetExecutor(),
	}
}

func (o *Orchestrator) Connect() {
	environmentTag := viper.GetString(flags.EnvironmentTag)
	connectionString := viper.GetString(flags.ConnectionString)
	isTechnicalAccount := viper.GetBool(flags.TechnicalAccount)
	isSecondary := viper.GetBool(flags.Secondary)
	printToStdout := viper.GetBool(flags.PrintToStdout)
	printConfigSilenceEverythingElse := viper.GetBool(flags.PrintToStdoutSilenceEverythingElse)

	validateConnectionString(connectionString, environmentTag, isTechnicalAccount)

	blob := downloadBlob(o, connectionString)
	validateDownloadedBlob(blob, isTechnicalAccount)

	configs := parseConfigs(blob)

	if !printConfigSilenceEverythingElse {
		logrus.Info("Using configuration last modified at " + configs.ModifiedAt)
	}

	configYaml := marshalToYaml(configs.PrimaryKubeConfig)

	if isSecondary {
		configYaml = marshalToYaml(configs.SecondaryKubeConfig)
	}

	if printToStdout || printConfigSilenceEverythingElse {
		if !printConfigSilenceEverythingElse {
			logrus.Info("===========================================================")
			logrus.Info("You can either replace your config file in $HOME/.kube/config manually, or merge the files.")
			logrus.Info("Check for reference: https://stackoverflow.com/questions/46184125/how-to-merge-kubectl-config-file-with-kube-config")
			logrus.Info("===========================================================")
			logrus.Info("==================== Kube Config:  ========================")
			logrus.Info("===========================================================")
		}

		fmt.Println(configYaml)
	} else {
		proceed := true
		if !viper.GetBool(flags.AutoApprove) {
			if !o.executor.AskUserToConfirm("Proceeding will overwrite your local $HOME/.kube/config file. Your old config will be backed up, but the impact of " +
				"this is that you will lose existing connections to other clusters. You can manually restore your connections by renaming the config " +
				"backup (in .kube directory) back to 'config' file name. Type 'yes' to proceed. You can also consider using the " +
				flags.PrintToStdout + "flag to see instructions on merging the kube config files") {
				proceed = false
			}
		}

		if proceed {
			saveKubeConfigToFile(configYaml)

			logrus.Info("Connection setup completed.")
		} else {
			logrus.Warn("Connection setup aborted.")
		}
	}
}

func downloadBlob(o *Orchestrator, connectionString string) string {
	blob, err := o.executor.ExecuteSilent("curl -s --retry 3 " + connectionString)

	if err != nil {
		panic(err)
	}

	return blob
}

func parseConfigs(yamlString string) KubeConfigsContainerV1 {
	var config KubeConfigsContainerV1

	err := yaml.Unmarshal([]byte(yamlString), &config)

	if err != nil {
		panic(err)
	}

	return config
}

func saveKubeConfigToFile(configYaml string) {
	home, err := homedir.Dir()
	panicOnError(err)

	err = os.MkdirAll(filepath.Join(home, ".kube"), os.ModePerm)
	panicOnError(err)

	configFilePath := filepath.Join(home, ".kube", "config")

	// check if file already there, make a backup
	if _, err := os.Stat(configFilePath); err == nil {
		copyFile(configFilePath, filepath.Join(home, ".kube", "copsctl_backup_config"))
	}

	err = os.WriteFile(configFilePath, []byte(configYaml), 0600)
	panicOnError(err)
}

func copyFile(sourceFile, destinationFile string) error {
	in, err := os.Open(sourceFile)

	if err != nil {
		return fmt.Errorf("could not open file which was expected to exist: %v, error: %w", sourceFile, err)
	}

	defer in.Close()

	out, err := os.Create(destinationFile)

	if err != nil {
		return fmt.Errorf("could not create the destination file: %v, error: %w", destinationFile, err)
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	if err != nil {
		return fmt.Errorf("problem copying the file from %v to %v, error: %w", sourceFile, destinationFile, err)
	}

	return out.Close()
}

type KubeConfigsContainerV1 struct {
	Version             string                 `yaml:"version"`
	ModifiedAt          string                 `yaml:"modifiedAt"`
	Type                string                 `yaml:"type"`
	PrimaryKubeConfig   map[string]interface{} `yaml:"primaryKubeConfig"`
	SecondaryKubeConfig map[string]interface{} `yaml:"secondaryKubeConfig"`
}
