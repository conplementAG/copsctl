package kubernetes

import (
	"encoding/json"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/conplementAG/copsctl/internal/common/commands"
	"github.com/conplementAG/copsctl/internal/common/file_processing"
	"github.com/conplementAG/copsctl/internal/common/logging"
)

func UseContext(contextName string) error {
	command := "kubectl config use-context " + contextName
	_, err := commands.ExecuteCommand(commands.Create(command))
	return err
}

func GetCurrentConfig() (*ConfigResponse, error) {
	command := "kubectl config view -o json"
	out, err := commands.ExecuteCommand(commands.Create(command))

	if err != nil {
		return nil, err
	}

	config := &ConfigResponse{}
	json.Unmarshal([]byte(out), &config)
	return config, nil
}

func PrintAllCopsNamespaces() error {
	command := "kubectl get copsnamespaces"
	out, err := commands.ExecuteCommand(commands.Create(command))

	if err != nil {
		return err
	}

	logging.Info("\nNamespaces:\n" + out)
	return nil
}

func GetCopsNamespace(namespace string) (*CopsNamespaceResponse, error) {
	command := "kubectl get CopsNamespace " + namespace + " -o json"
	out, err := commands.ExecuteCommand(commands.Create(command))

	if err != nil {
		return nil, err
	}

	response := &CopsNamespaceResponse{}
	json.Unmarshal([]byte(out), &response)
	return response, nil
}

func Apply(filepath string) (string, error) {
	command := "kubectl apply -f " + filepath
	return commands.ExecuteCommandLongRunning(commands.Create(command))
}

func Delete(filepath string) (string, error) {
	command := "kubectl delete --wait -f " + filepath
	return commands.ExecuteCommandLongRunning(commands.Create(command))
}

func ApplyString(content string) (string, error) {
	temporaryDirectory, temporaryFile := file_processing.WriteStringToTemporaryFile(content, "resource.yaml")
	defer file_processing.DeletePath(temporaryDirectory)

	return Apply(temporaryFile)
}

func DeleteString(content string) (string, error) {
	temporaryDirectory, temporaryFile := file_processing.WriteStringToTemporaryFile(content, "resource.yaml")
	defer file_processing.DeletePath(temporaryDirectory)

	return Delete(temporaryFile)
}

func CanIGetPods(namespace string) bool {
	data, err := commands.ExecuteCommandLongRunning(commands.Create("kubectl auth can-i get pods -n " + namespace))
	return err == nil && strings.TrimSuffix(data, "\n") == "yes"
}

func GetCurrentMasterPlaneFqdn() (string, error) {
	currentConfig, err := GetCurrentConfig()

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

func CreateServiceAccount(namespace string, accountName string) error {
	command := "kubectl create serviceaccount " + accountName + " --namespace " + namespace
	_, err := commands.ExecuteCommand(commands.Create(command))
	return err
}

func RemoveServiceAccount(namespace string, accountName string) error {
	command := "kubectl delete serviceaccount " + accountName + " --namespace " + namespace
	_, err := commands.ExecuteCommand(commands.Create(command))
	return err
}

func DeleteDeployment(deploymentName string, namespace string) error {
	command := "kubectl delete deployment " + deploymentName + " -n " + namespace
	_, err := commands.ExecuteCommand(commands.Create(command))
	return err
}
