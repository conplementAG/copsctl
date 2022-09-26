package kubernetes

import (
	"encoding/json"
	"github.com/conplementag/cops-hq/pkg/commands"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/conplementAG/copsctl/internal/common/file_processing"
)

func GetCurrentConfig(executor commands.Executor) (*ConfigResponse, error) {
	command := "kubectl config view -o json"
	out, err := executor.Execute(command)

	if err != nil {
		return nil, err
	}

	config := &ConfigResponse{}
	json.Unmarshal([]byte(out), &config)
	return config, nil
}

func PrintAllCopsNamespaces(executor commands.Executor) error {
	command := "kubectl get copsnamespaces"
	out, err := executor.Execute(command)

	if err != nil {
		return err
	}

	logrus.Info("\nNamespaces:\n" + out)
	return nil
}

func GetCopsNamespace(executor commands.Executor, namespace string) (*CopsNamespaceResponse, error) {
	command := "kubectl get CopsNamespace " + namespace + " -o json"
	out, err := executor.Execute(command)

	if err != nil {
		return nil, err
	}

	response := &CopsNamespaceResponse{}
	json.Unmarshal([]byte(out), &response)
	return response, nil
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
	defer file_processing.DeletePath(temporaryDirectory)

	return Apply(executor, temporaryFile)
}

func DeleteString(executor commands.Executor, content string) (string, error) {
	temporaryDirectory, temporaryFile := file_processing.WriteStringToTemporaryFile(content, "resource.yaml")
	defer file_processing.DeletePath(temporaryDirectory)

	return Delete(executor, temporaryFile)
}

func CanIGetPods(executor commands.Executor, namespace string) bool {
	data, err := executor.ExecuteWithProgressInfo("kubectl auth can-i get pods -n " + namespace)
	return err == nil && strings.TrimSuffix(data, "\n") == "yes"
}

func GetCurrentMasterPlaneFqdn(executor commands.Executor) (string, error) {
	currentConfig, err := GetCurrentConfig(executor)

	if err != nil {
		return "", err
	}

	currentContextName := currentConfig.CurrentContext
	currentContextResponse := linq.From(currentConfig.Contexts).SingleWithT(func(c Context) bool {
		return c.Name == currentContextName
	}).(Context)
	currentClusterResponse := linq.From(currentConfig.Clusters).SingleWithT(func(c Cluster) bool {
		return c.Name == currentContextResponse.Context.Cluster
	}).(Cluster)

	return currentClusterResponse.Cluster.Server, nil
}
