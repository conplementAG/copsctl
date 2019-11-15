package kubernetes

import (
	"github.com/conplementAG/copsctl/internal/common/commands"
	yaml "gopkg.in/yaml.v2"
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

func GetServiceAccountSecret(namespace string, name string) ServiceAccountSecret {
	out, err := commands.ExecuteCommand(commands.Createf("kubectl get secret %s -n %s -o yaml --ignore-not-found", name, namespace))

	if err != nil {
		panic("we should never get an err when using ignore-not-found " + err.Error())
	}

	if out == "" {
		return ServiceAccountSecret{}
	}

	var secret ServiceAccountSecret
	yaml.Unmarshal([]byte(out), &secret)

	return secret
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

type ServiceAccountSecret struct {
	APIVersion string `yaml:"apiVersion"`
	Data       struct {
		CaCrt     string `yaml:"ca.crt"`
		Namespace string `yaml:"namespace"`
		Token     string `yaml:"token"`
	} `yaml:"data"`
	Kind     string `yaml:"kind"`
	Metadata struct {
		Annotations struct {
			KubernetesIoServiceAccountName string `yaml:"kubernetes.io/service-account.name"`
			KubernetesIoServiceAccountUID  string `yaml:"kubernetes.io/service-account.uid"`
		} `yaml:"annotations"`
		CreationTimestamp struct {
		} `yaml:"creationTimestamp"`
		Name            string `yaml:"name"`
		Namespace       string `yaml:"namespace"`
		ResourceVersion string `yaml:"resourceVersion"`
		SelfLink        string `yaml:"selfLink"`
		UID             string `yaml:"uid"`
	} `yaml:"metadata"`
	Type string `yaml:"type"`
}
