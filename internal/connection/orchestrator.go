package connection

import (
	"fmt"
	"github.com/conplementAG/copsctl/internal/cmd/flags"
	"github.com/conplementAG/copsctl/internal/common/commands"
	"github.com/conplementAG/copsctl/internal/common/logging"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Connect() {
	environmentTag := viper.GetString(flags.EnvironmentTag)
	connectionString := viper.GetString(flags.ConnectionString)
	isTechnicalAccount := viper.GetBool(flags.TechnicalAccount)
	isSecondary := viper.GetBool(flags.Secondary)
	printToStdout := viper.GetBool(flags.PrintToStdout)
	printConfigSilenceEverythingElse := viper.GetBool(flags.PrintToStdoutSilenceEverythingElse)

	validateConnectionString(connectionString, environmentTag, isTechnicalAccount)

	blob := downloadBlob(connectionString)
	validateDownloadedBlob(blob, isTechnicalAccount)

	configs := parseConfigs(blob)

	if !printConfigSilenceEverythingElse {
		logging.Info("Using configuration last modified at " + configs.ModifiedAt)
	}

	configYaml := marshalToYaml(configs.PrimaryKubeConfig)

	if isSecondary {
		configYaml = marshalToYaml(configs.SecondaryKubeConfig)
	}

	if printToStdout || printConfigSilenceEverythingElse {
		if !printConfigSilenceEverythingElse {
			logging.Info("===========================================================")
			logging.Info("You can either replace your config file in $HOME/.kube/config manually, or merge the files.")
			logging.Info("Check for reference: https://stackoverflow.com/questions/46184125/how-to-merge-kubectl-config-file-with-kube-config")
			logging.Info("===========================================================")
			logging.Info("==================== Kube Config:  ========================")
			logging.Info("===========================================================")
		}

		fmt.Println(configYaml)
	} else {
		if !viper.GetBool(flags.AutoApprove) {
			confirmOperation("Proceeding will overwrite your local $HOME/.kube/config file. Your old config will be backed up, but the impact of " +
				"this is that you will lose existing connections to other clusters. You can manually restore your connections by renaming the config " +
				"backup (in .kube directory) back to 'config' file name. Type 'yes' to proceed. You can also consider using the " +
				flags.PrintToStdout + "flag to see instructions on merging the kube config files")
		}

		saveKubeConfigToFile(configYaml)

		logging.Info("Connection setup completed.")
	}
}

func downloadBlob(connectionString string) string {
	blob, err := commands.ExecuteCommandWithSecretContents(
		commands.Create("curl -s --retry 3 " + connectionString))

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

	err = ioutil.WriteFile(configFilePath, []byte(configYaml), 0600)
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
