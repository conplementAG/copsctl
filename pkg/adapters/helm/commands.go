package helm

import (
	"github.com/conplementAG/copsctl/pkg/common/commands"
)

func InitHelm(serviceAccount string, namespace string) {
	command := "helm init --service-account " + serviceAccount + " --tiller-namespace " + namespace
	commands.ExecuteCommand(commands.Create(command))
}
