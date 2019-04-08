package kubernetes

import (
	"encoding/json"

	"github.com/conplementAG/copsctl/pkg/common/commands"
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

func Apply(filepath string) string {
	command := "kubectl apply -f " + filepath
	data := commands.ExecuteCommandLongRunning(commands.Create(command))
	return data
}
