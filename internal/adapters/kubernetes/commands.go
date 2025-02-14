package kubernetes

import (
	"encoding/json"
	"github.com/conplementAG/copsctl/internal/common"
	"github.com/conplementAG/copsctl/internal/common/file_processing"
	"github.com/conplementag/cops-hq/v2/pkg/commands"
	"github.com/sirupsen/logrus"
	"strings"
)

func PrintAllCopsNamespaces(executor commands.Executor) error {
	command := "kubectl get cns"
	out, err := executor.Execute(command)

	if err != nil {
		return err
	}

	logrus.Info("\nNamespaces:\n" + out)
	return nil
}

func ExistsCopsNamespace(executor commands.Executor, namespace string) (bool, error) {
	response, err := GetCopsNamespace(executor, namespace)

	if err != nil {
		return false, err
	}

	if response != nil {
		return true, nil
	} else {
		return false, nil
	}
}

func GetCopsNamespace(executor commands.Executor, namespace string) (*CopsNamespaceResponse, error) {
	// ignore-not-found flag needed to avoid command to fail with hard panic mode (because of globally set PanicOnAnyError=true)
	command := "kubectl get cns " + namespace + " -o json --ignore-not-found"
	out, err := executor.Execute(command)

	if err != nil {
		return nil, err
	}

	out = strings.TrimSuffix(out, "\n")

	if out == "" {
		return nil, nil
	} else {
		response := &CopsNamespaceResponse{}
		err := json.Unmarshal([]byte(out), &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

func Apply(executor commands.Executor, filepath string) (string, error) {
	command := "kubectl apply -f " + filepath
	return executor.ExecuteWithProgressInfo(command)
}

func Delete(executor commands.Executor, filepath string) (string, error) {
	command := "kubectl delete --wait -f " + filepath
	return executor.ExecuteWithProgressInfo(command)
}

func ApplyString(executor commands.Executor, content string) (string, error) {
	temporaryDirectory, temporaryFile := file_processing.WriteStringToTemporaryFile(content, "resource.yaml")
	defer func() {
		err := file_processing.DeletePath(temporaryDirectory)
		common.FatalOnError(err)
	}()

	return Apply(executor, temporaryFile)
}

func DeleteString(executor commands.Executor, content string) (string, error) {
	temporaryDirectory, temporaryFile := file_processing.WriteStringToTemporaryFile(content, "resource.yaml")
	defer func() {
		err := file_processing.DeletePath(temporaryDirectory)
		common.FatalOnError(err)
	}()

	return Delete(executor, temporaryFile)
}

func CanIGetPods(executor commands.Executor, namespace string) bool {
	data, err := executor.ExecuteWithProgressInfo("kubectl auth can-i get pods -n " + namespace)
	return err == nil && strings.TrimSuffix(data, "\n") == "yes"
}
