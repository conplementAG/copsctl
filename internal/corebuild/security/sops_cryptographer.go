package security

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type sops struct {
	configFilePath string
}

func NewSopsCryptographer(sopsConfigFilePath string) Cryptographer {
	return &sops{
		configFilePath: sopsConfigFilePath,
	}
}

func (c *sops) DecryptYamlContent(cipherTextYaml string) (string, error) {
	tempFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		return "", err
	}

	tempFilePath := tempFile.Name()
	defer func() {
		err = os.RemoveAll(tempFilePath)
	}()

	if _, err := tempFile.WriteString(cipherTextYaml); err != nil {
		return "", err
	}
	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	// Decrypt the file using SOPS
	cmd := exec.Command("sops", "--decrypt", "--config", c.configFilePath, tempFilePath)
	decryptedData, err := cmd.CombinedOutput()
	if strings.HasPrefix(string(decryptedData), "sops metadata not found") {
		logrus.Warningf("Ciphertext was not sops encryted. Ciphertext will be returned.")
		decryptedData = []byte(cipherTextYaml)
	} else if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}
