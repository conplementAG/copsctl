package helm

import (
	"github.com/conplementAG/copsctl/internal/common/commands"
)

func InitHelm(serviceAccount string, namespace string) (string, error) {
	command := "helm init --wait --upgrade --service-account " + serviceAccount + " --tiller-namespace " + namespace
	return commands.ExecuteCommandLongRunning(commands.Create(command))
}
