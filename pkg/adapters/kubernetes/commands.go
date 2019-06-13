package kubernetes

import (
	"encoding/json"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/conplementAG/copsctl/pkg/common/commands"
	"github.com/conplementAG/copsctl/pkg/common/fileprocessing"
)

// UseContext sets the given context as the current context in the config file
func UseContext(contextName string) {
	command := "kubectl config use-context " + contextName
	commands.ExecuteCommand(commands.Create(command))
}

// GetCurrentConfig gets the current config
func GetCurrentConfig() *ConfigResponse {
	command := "kubectl config view -o json"
	out := commands.ExecuteCommand(commands.Create(command))

	config := &ConfigResponse{}
	json.Unmarshal([]byte(out), &config)
	return config
}

// GetCopsNamespace gets the given CopsNamespace
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

func ApplyString(content string) {
	temporaryDirectory, temporaryFile := fileprocessing.WriteStringToTemporaryFile(content, "resource.yaml")
	Apply(temporaryFile)
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
