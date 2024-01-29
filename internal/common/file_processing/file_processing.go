package file_processing

import (
	"embed"
	"fmt"
	"github.com/conplementAG/copsctl/internal/common"
	"github.com/conplementAG/copsctl/internal/corebuild/security"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
)

// WriteStringToTemporaryFile writes the file contents into a file on a temporary disk location
func WriteStringToTemporaryFile(fileContents string, filePath string) (outputFolder string, outputFile string) {
	outputFolderPath := createUniqueDirectory()

	generatedFilePath := filepath.Join(outputFolderPath, filePath)
	err := os.WriteFile(generatedFilePath, []byte(fileContents), 0644)
	common.FatalOnError(err)

	return outputFolderPath, generatedFilePath
}

// DeletePath deletes the file from the disk
func DeletePath(filePath string) {
	err := os.RemoveAll(filePath)
	common.FatalOnError(err)
}

// InterpolateStaticFiles loads all the files in given embed FS path.
// It depends on resource embedding, set by go:embed directive
// Replaces the variables based on the given dictionary,
// and returns the path to the generated directory where the results are stored
func InterpolateStaticFiles(inputPathFs embed.FS, inputPathRootFolderName string, variables map[string]string) string {
	directory, readDirError := inputPathFs.ReadDir(inputPathRootFolderName)
	common.FatalOnError(readDirError)

	uniqueOutputFolder := createUniqueDirectory()

	for _, file := range directory {
		f, err := inputPathFs.Open(inputPathRootFolderName + "/" + file.Name())
		common.FatalOnError(err)
		filesContent, err := io.ReadAll(f)
		common.FatalOnError(err)
		fileContentString := string(filesContent)
		for key, value := range variables {
			fileContentString = strings.Replace(fileContentString, key, value, -1)
		}

		err = os.WriteFile(filepath.Join(uniqueOutputFolder, file.Name()), []byte(fileContentString), 0644)
		common.FatalOnError(err)
	}

	return uniqueOutputFolder
}

func LoadEncryptedFile[T interface{}](filename string, cryptographer security.Cryptographer) (*T, error) {
	encryptedContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content, err := cryptographer.DecryptYamlContent(string(encryptedContent))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content %s: %w", string(encryptedContent), err)
	}

	var result *T
	err = yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return result, nil
}

func createUniqueDirectory() string {
	folderPath := filepath.Join(".", ".generated", xid.New().String())
	err := os.MkdirAll(folderPath, os.ModePerm)
	common.FatalOnError(err)
	return folderPath
}
