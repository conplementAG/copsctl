package kubernetes

import (
	"github.com/conplementAG/copsctl/internal/common/commands"
	"gopkg.in/yaml.v2"
)

func GetServiceAccount(namespace string, name string) ServiceAccount {
	out, err := commands.ExecuteCommand(commands.Createf("kubectl get serviceaccount %s -n %s -o yaml --ignore-not-found", name, namespace))

	if err != nil {
		panic("we should never get an err when using ignore-not-found " + err.Error())
	}

	if out == "" {
		return ServiceAccount{}
	}

	var serviceAccount ServiceAccount
	yaml.Unmarshal([]byte(out), &serviceAccount)

	return serviceAccount
}

type ServiceAccount struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		CreationTimestamp struct {
		} `yaml:"creationTimestamp"`
		Name            string `yaml:"name"`
		Namespace       string `yaml:"namespace"`
		ResourceVersion string `yaml:"resourceVersion"`
		SelfLink        string `yaml:"selfLink"`
		UID             string `yaml:"uid"`
	} `yaml:"metadata"`
	Secrets []struct {
		Name string `yaml:"name"`
	} `yaml:"secrets"`
}
