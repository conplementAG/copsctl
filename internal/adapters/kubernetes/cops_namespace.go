package kubernetes

type CopsNamespaceResponse struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		NamespaceAdminUsers           []string             `yaml:"namespaceAdminUsers"`
		NamespaceAdminServiceAccounts []CopsServiceAccount `yaml:"namespaceAdminServiceAccounts"`
	} `yaml:"spec"`
}

type CopsServiceAccount struct {
	ServiceAccount string `yaml:"serviceAccount"`
	Namespace      string `yaml:"namespace"`
}

func (o CopsNamespaceResponse) GetProjectName() string {
	return o.Metadata.Labels["conplement.de/projectName"]
}

func (o CopsNamespaceResponse) GetProjectCostCenter() string {
	return o.Metadata.Labels["conplement.de/projectCostCenter"]
}
