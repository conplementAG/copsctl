package file_processing

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
)

// WriteStringToTemporaryFile writes the file contents into a file on a temporary disk location
func WriteStringToTemporaryFile(fileContents string, filePath string) (outputfolder string, outputfile string) {
	outputFolderPath := createUniqueDirectory()

	generatedFilePath := filepath.Join(outputFolderPath, filePath)
	err := ioutil.WriteFile(generatedFilePath, []byte(fileContents), 0644)
	panicOnError(err)

	return outputFolderPath, generatedFilePath
}

// DeletePath deletes the file from the disk
func DeletePath(filePath string) {
	err := os.RemoveAll(filePath)
	panicOnError(err)
}

// InterpolateStaticFiles loads all the files in given resource path.
// E.g.: internal/azure_devops/global
// It depends on resource embedding, that can be triggered by go generate.
// Replaces the variables based on the given dictionary,
// and returns the path to the generated directory where the results are stored
func InterpolateStaticFiles(inputPathFs embed.FS, inputPath string, variables map[string]string) string {
	directory, readDirError := inputPathFs.ReadDir(inputPath)
	panicOnError(readDirError)

	uniqueOutputFolder := createUniqueDirectory()

	for _, file := range directory {
		f, erri := inputPathFs.Open(inputPath + "/" + file.Name())
		if erri != nil {
			panicOnError(erri)
		}
		filesContent, _ := ioutil.ReadAll(f)
		fileContentString := string(filesContent)
		for key, value := range variables {
			fileContentString = strings.Replace(fileContentString, key, value, -1)
		}

		err := ioutil.WriteFile(filepath.Join(uniqueOutputFolder, file.Name()), []byte(fileContentString), 0644)
		panicOnError(err)
	}

	return uniqueOutputFolder
}

func panicOnError(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func createUniqueDirectory() string {
	folderPath := filepath.Join(".", ".generated", xid.New().String())
	err := os.MkdirAll(folderPath, os.ModePerm)
	panicOnError(err)
	return folderPath
}
