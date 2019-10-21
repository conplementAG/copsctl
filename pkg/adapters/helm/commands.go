package helm

import (
	"github.com/conplementAG/copsctl/pkg/common/commands"
)

func InitHelm(serviceAccount string, namespace string) (string, error) {
	command := "helm init --service-account " + serviceAccount + " --tiller-namespace " + namespace
	return commands.ExecuteCommand(commands.Create(command))
}
