package file_processing

import (
	"embed"
	"fmt"
	"github.com/conplementAG/copsctl/internal/common"
	"github.com/conplementAG/copsctl/internal/corebuild/security"
	"github.com/rs/xid"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
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
func DeletePath(filePath string) error {
	return os.RemoveAll(filePath)
}

// CreateTempDirectory loads all the files in given embed FS path.
// It depends on resource embedding, set by go:embed directive
// Returns the path to the generated directory where the results are stored
func CreateTempDirectory(inputPathFs embed.FS, inputPathRootFolderName string) (string, error) {

	uniqueOutputFolder := createUniqueDirectory()

	err := fs.WalkDir(inputPathFs, inputPathRootFolderName, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(inputPathRootFolderName, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(uniqueOutputFolder, relPath)
		if entry.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}
		fileContent, err := inputPathFs.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, fileContent, 0644)
	})

	return uniqueOutputFolder, err
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
