package kubernetes

import (
	"fmt"
	"github.com/conplementag/cops-hq/pkg/commands"
	"gopkg.in/yaml.v2"
)

func GetServiceAccountSecret(executor commands.Executor, namespace string, name string) ServiceAccountSecret {
	out, err := executor.Execute(fmt.Sprintf("kubectl get secret %s -n %s -o yaml --ignore-not-found", name, namespace))

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
