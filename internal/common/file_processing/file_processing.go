package file_processing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/conplementAG/copsctl/internal/resources"
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
func InterpolateStaticFiles(inputPath string, variables map[string]string) string {
	filesystem := resources.Dir(false, "/")
	directory, openDirError := filesystem.Open(inputPath)
	panicOnError(openDirError)

	files, listDirectoryError := directory.Readdir(9999)
	panicOnError(listDirectoryError)

	uniqueOutputFolder := createUniqueDirectory()

	for _, f := range files {

		fileContents, _ := resources.FSString(false, "/"+inputPath+"/"+f.Name())

		for key, value := range variables {
			fileContents = strings.Replace(fileContents, key, value, -1)
		}

		err := ioutil.WriteFile(filepath.Join(uniqueOutputFolder, f.Name()), []byte(fileContents), 0644)
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
