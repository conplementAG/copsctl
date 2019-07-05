package kubernetes

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/conplementAG/copsctl/pkg/common/commands"
	"github.com/conplementAG/copsctl/pkg/common/fileprocessing"
)

func UseContext(contextName string) {
	command := "kubectl config use-context " + contextName
	commands.ExecuteCommand(commands.Create(command))
}

func GetCurrentConfig() *ConfigResponse {
	command := "kubectl config view -o json"
	out := commands.ExecuteCommand(commands.Create(command))

	config := &ConfigResponse{}
	json.Unmarshal([]byte(out), &config)
	return config
}

func PrintAllCopsNamespaces() {
	command := "kubectl get copsnamespaces"
	out := commands.ExecuteCommand(commands.Create(command))
	log.Println(out)
}

func GetCopsNamespace(namespace string) *CopsNamespaceResponse {
	command := "kubectl get CopsNamespace " + namespace + " -o json"
	out := commands.ExecuteCommand(commands.Create(command))
	response := &CopsNamespaceResponse{}
	json.Unmarshal([]byte(out), &response)
	return response
}

func Apply(filepath string) string {
	command := "kubectl apply -f " + filepath
	data := commands.ExecuteCommandLongRunning(commands.Create(command))
	return data
}

func Delete(filepath string) string {
	command := "kubectl delete -f " + filepath
	data := commands.ExecuteCommandLongRunning(commands.Create(command))
	return data
}

func ApplyString(content string) {
	temporaryDirectory, temporaryFile := fileprocessing.WriteStringToTemporaryFile(content, "resource.yaml")
	Apply(temporaryFile)
	fileprocessing.DeletePath(temporaryDirectory)
}

func DeleteString(content string) {
	temporaryDirectory, temporaryFile := fileprocessing.WriteStringToTemporaryFile(content, "resource.yaml")
	Delete(temporaryFile)
	fileprocessing.DeletePath(temporaryDirectory)
}

func CanIGetPods(namespace string) bool {
	data := commands.ExecuteCommand(commands.Create("kubectl auth can-i get pods -n " + namespace))
	return strings.TrimSuffix(data, "\n") == "yes"
}

func GetCurrentMasterPlaneFqdn() string {
	currentConfig := GetCurrentConfig()
	currentContextName := currentConfig.CurrentContext
	currentContextResponse := linq.From(currentConfig.Contexts).SingleWithT(func(c Context) bool {
		return c.Name == currentContextName
	}).(Context)
	currentClusterResponse := linq.From(currentConfig.Clusters).SingleWithT(func(c Cluster) bool {
		return c.Name == currentContextResponse.Context.Cluster
	}).(Cluster)

	return currentClusterResponse.Cluster.Server
}

func CreateServiceAccount(namespace string, accountName string) {
	command := "kubectl create serviceaccount " + accountName + " --namespace " + namespace
	commands.ExecuteCommand(commands.Create(command))
}

func RemoveServiceAccount(namespace string, accountName string) {
	command := "kubectl delete serviceaccount " + accountName + " --namespace " + namespace
	commands.ExecuteCommand(commands.Create(command))
}

func DeleteDeployment(deploymentName string, namespace string) {
	command := "kubectl delete deployment " + deploymentName + " -n " + namespace
	commands.ExecuteCommand(commands.Create(command))
}
